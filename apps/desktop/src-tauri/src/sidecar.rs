use std::time::Duration;

use tauri::{AppHandle, Emitter, Manager, State};
use tauri_plugin_shell::process::{CommandChild, CommandEvent};
use tauri_plugin_shell::ShellExt;
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::sync::Mutex;
use tokio::time::{sleep, timeout};

#[derive(serde::Serialize, Clone, Default)]
#[serde(rename_all = "camelCase")]
pub struct SidecarStatus {
    pub running: bool,
    pub healthy: bool,
    pub restarts: u32,
    pub pid: Option<u32>,
    pub last_error: Option<String>,
    pub port: u16,
}

struct SidecarInner {
    status: SidecarStatus,
    child: Option<CommandChild>,
    stop_requested: bool,
}

pub struct SidecarState {
    inner: Mutex<SidecarInner>,
}

impl SidecarState {
    pub fn new() -> Self {
        Self {
            inner: Mutex::new(SidecarInner {
                status: SidecarStatus {
                    port: 8080,
                    ..Default::default()
                },
                child: None,
                stop_requested: false,
            }),
        }
    }
}

pub fn supervise(app: AppHandle) {
    tauri::async_runtime::spawn(async move {
        let mut backoff = Duration::from_secs(1);
        loop {
            let stopped = {
                let state: State<'_, SidecarState> = app.state();
                let stopped = state.inner.lock().await.stop_requested;
                stopped
            };
            if stopped {
                sleep(Duration::from_secs(2)).await;
                continue;
            }
            match run_sidecar(&app).await {
                Ok(()) => {
                    backoff = Duration::from_secs(1);
                }
                Err(message) => {
                    set_status(&app, |status| {
                        status.running = false;
                        status.healthy = false;
                        status.last_error = Some(message);
                    })
                    .await;
                }
            }
            {
                let state: State<'_, SidecarState> = app.state();
                let mut inner = state.inner.lock().await;
                inner.child = None;
                inner.status.running = false;
                inner.status.healthy = false;
                inner.status.pid = None;
                inner.status.restarts += 1;
                let snapshot = inner.status.clone();
                drop(inner);
                let _ = app.emit("sidecar-status", snapshot);
            }
            sleep(backoff).await;
            backoff = std::cmp::min(backoff * 2, Duration::from_secs(30));
        }
    });
}

async fn set_status(app: &AppHandle, mutate: impl FnOnce(&mut SidecarStatus)) {
    let state: State<'_, SidecarState> = app.state();
    let mut inner = state.inner.lock().await;
    mutate(&mut inner.status);
    let snapshot = inner.status.clone();
    drop(inner);
    let _ = app.emit("sidecar-status", snapshot);
}

async fn run_sidecar(app: &AppHandle) -> Result<(), String> {
    let port = {
        let state: State<'_, SidecarState> = app.state();
        let port = state.inner.lock().await.status.port;
        port
    };
    let (mut rx, child) = app
        .shell()
        .sidecar("codedockd")
        .map_err(|err| format!("sidecar unavailable: {err}"))?
        .spawn()
        .map_err(|err| format!("spawn daemon: {err}"))?;
    let pid = child.pid();
    {
        let state: State<'_, SidecarState> = app.state();
        let mut inner = state.inner.lock().await;
        inner.child = Some(child);
        inner.status.running = true;
        inner.status.last_error = None;
        inner.status.pid = Some(pid);
        let snapshot = inner.status.clone();
        drop(inner);
        let _ = app.emit("sidecar-status", snapshot);
    }
    loop {
        let healthy = probe_health(port).await;
        set_status(app, |status| {
            status.healthy = healthy;
            if healthy {
                status.last_error = None;
            }
        })
        .await;
        tokio::select! {
            event = rx.recv() => {
                match event {
                    Some(CommandEvent::Terminated(_)) => return Ok(()),
                    Some(CommandEvent::Error(message)) => return Err(message),
                    Some(_) => {}
                    None => return Ok(()),
                }
            }
            _ = sleep(Duration::from_secs(2)) => {}
        }
    }
}

async fn probe_health(port: u16) -> bool {
    let dial = timeout(
        Duration::from_secs(2),
        tokio::net::TcpStream::connect(format!("127.0.0.1:{port}")),
    )
    .await;
    let mut stream = match dial {
        Ok(Ok(stream)) => stream,
        _ => return false,
    };
    if stream
        .write_all(b"GET /healthz HTTP/1.0\r\nHost: localhost\r\n\r\n")
        .await
        .is_err()
    {
        return false;
    }
    let mut head = [0u8; 15];
    match timeout(Duration::from_secs(2), stream.read_exact(&mut head)).await {
        Ok(Ok(_)) => head.starts_with(b"HTTP/1.0 200") || head.starts_with(b"HTTP/1.1 200"),
        _ => false,
    }
}

#[tauri::command]
pub async fn sidecar_status(state: State<'_, SidecarState>) -> Result<SidecarStatus, String> {
    Ok(state.inner.lock().await.status.clone())
}

#[tauri::command]
pub async fn sidecar_restart(
    app: AppHandle,
    state: State<'_, SidecarState>,
) -> Result<(), String> {
    let mut inner = state.inner.lock().await;
    if let Some(child) = inner.child.take() {
        child.kill().map_err(|err| err.to_string())?;
    }
    inner.status.last_error = None;
    let snapshot = inner.status.clone();
    drop(inner);
    let _ = app.emit("sidecar-status", snapshot);
    Ok(())
}

#[tauri::command]
pub async fn sidecar_stop(state: State<'_, SidecarState>) -> Result<(), String> {
    let mut inner = state.inner.lock().await;
    inner.stop_requested = true;
    if let Some(child) = inner.child.take() {
        child.kill().map_err(|err| err.to_string())?;
    }
    inner.status.running = false;
    inner.status.healthy = false;
    Ok(())
}

#[tauri::command]
pub async fn sidecar_start(state: State<'_, SidecarState>) -> Result<(), String> {
    state.inner.lock().await.stop_requested = false;
    Ok(())
}

#[tauri::command]
pub async fn sidecar_set_port(state: State<'_, SidecarState>, port: u16) -> Result<(), String> {
    if port == 0 {
        return Err("port is required".to_string());
    }
    let mut inner = state.inner.lock().await;
    inner.status.port = port;
    if let Some(child) = inner.child.take() {
        child.kill().map_err(|err| err.to_string())?;
    }
    Ok(())
}
