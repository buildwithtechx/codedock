use tauri::AppHandle;
use tauri_plugin_store::StoreExt;

const STORE_PATH: &str = "connection.json";
const KEYRING_SERVICE: &str = "run.codedock.desktop";
const KEYRING_USER: &str = "daemon-token";

#[derive(serde::Serialize, Clone, Default)]
#[serde(rename_all = "camelCase")]
pub struct DaemonConnection {
    pub daemon_url: String,
    pub daemon_port: u16,
    pub has_token: bool,
}

fn store(app: &AppHandle) -> Result<std::sync::Arc<tauri_plugin_store::Store<tauri::Wry>>, String> {
    app.store(STORE_PATH).map_err(|err| err.to_string())
}

#[tauri::command]
pub fn get_connection(app: AppHandle) -> Result<DaemonConnection, String> {
    let store = store(&app)?;
    let url = store
        .get("daemonUrl")
        .and_then(|value| value.as_str().map(str::to_string))
        .unwrap_or_else(|| "http://localhost:8080".to_string());
    let port = store
        .get("daemonPort")
        .and_then(|value| value.as_u64())
        .unwrap_or(8080) as u16;
    let entry = keyring::Entry::new(KEYRING_SERVICE, KEYRING_USER).map_err(|err| err.to_string())?;
    Ok(DaemonConnection {
        daemon_url: url,
        daemon_port: port,
        has_token: entry.get_password().is_ok(),
    })
}

#[tauri::command]
pub fn set_daemon_url(app: AppHandle, url: String) -> Result<(), String> {
    let trimmed = url.trim().trim_end_matches('/').to_string();
    if trimmed.is_empty() {
        return Err("daemon url is required".to_string());
    }
    if !trimmed.starts_with("http://") && !trimmed.starts_with("https://") {
        return Err("daemon url must start with http(s)://".to_string());
    }
    let store = store(&app)?;
    store.set("daemonUrl", serde_json::json!(trimmed));
    store.save().map_err(|err| err.to_string())?;
    Ok(())
}

#[tauri::command]
pub fn set_daemon_port(app: AppHandle, port: u16) -> Result<(), String> {
    if port == 0 {
        return Err("port is required".to_string());
    }
    let store = store(&app)?;
    store.set("daemonPort", serde_json::json!(port));
    store.save().map_err(|err| err.to_string())?;
    Ok(())
}

#[tauri::command]
pub fn save_token(token: String) -> Result<(), String> {
    if token.is_empty() {
        return Err("token is required".to_string());
    }
    let entry = keyring::Entry::new(KEYRING_SERVICE, KEYRING_USER).map_err(|err| err.to_string())?;
    entry
        .set_password(&token)
        .map_err(|err| err.to_string())?;
    Ok(())
}

#[tauri::command]
pub fn load_token() -> Result<String, String> {
    let entry = keyring::Entry::new(KEYRING_SERVICE, KEYRING_USER).map_err(|err| err.to_string())?;
    entry.get_password().map_err(|err| err.to_string())
}

#[tauri::command]
pub fn clear_token() -> Result<(), String> {
    let entry = keyring::Entry::new(KEYRING_SERVICE, KEYRING_USER).map_err(|err| err.to_string())?;
    match entry.delete_credential() {
        Ok(()) => Ok(()),
        Err(keyring::Error::NoEntry) => Ok(()),
        Err(err) => Err(err.to_string()),
    }
}
