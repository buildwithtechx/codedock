use tauri::{AppHandle, Emitter};
use tauri_plugin_updater::UpdaterExt;

#[derive(serde::Serialize, Clone, Default)]
#[serde(rename_all = "camelCase")]
pub struct UpdateInfo {
    pub available: bool,
    pub version: String,
    pub notes: String,
    pub date: String,
}

#[tauri::command]
pub async fn check_update(app: AppHandle) -> Result<UpdateInfo, String> {
    let updater = app.updater().map_err(|err| err.to_string())?;
    let update = updater.check().await.map_err(|err| err.to_string())?;
    match update {
        Some(update) => Ok(UpdateInfo {
            available: true,
            version: update.version.clone(),
            notes: update.body.clone().unwrap_or_default(),
            date: update.date.map(|d| d.to_string()).unwrap_or_default(),
        }),
        None => Ok(UpdateInfo::default()),
    }
}

#[tauri::command]
pub async fn download_and_install(app: AppHandle) -> Result<(), String> {
    let updater = app.updater().map_err(|err| err.to_string())?;
    let update = updater
        .check()
        .await
        .map_err(|err| err.to_string())?
        .ok_or_else(|| "no update available".to_string())?;
    let handle = app.clone();
    update
        .download_and_install(
            move |chunk_length, content_length| {
                let _ = handle.emit(
                    "update-status",
                    serde_json::json!({
                        "state": "downloading",
                        "chunk": chunk_length,
                        "total": content_length,
                    }),
                );
            },
            || {},
        )
        .await
        .map_err(|err| err.to_string())?;
    let _ = app.emit(
        "update-status",
        serde_json::json!({ "state": "ready" }),
    );
    Ok(())
}

#[tauri::command]
pub fn restart_app(app: AppHandle) {
    app.restart();
}
