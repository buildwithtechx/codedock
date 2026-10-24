#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod connection;
mod sidecar;
mod updates;

fn main() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_updater::Builder::new().build())
        .plugin(tauri_plugin_store::Builder::new().build())
        .manage(sidecar::SidecarState::new())
        .invoke_handler(tauri::generate_handler![
            sidecar::sidecar_status,
            sidecar::sidecar_restart,
            sidecar::sidecar_stop,
            sidecar::sidecar_start,
            sidecar::sidecar_set_port,
            connection::get_connection,
            connection::set_daemon_url,
            connection::set_daemon_port,
            connection::save_token,
            connection::load_token,
            connection::clear_token,
            updates::check_update,
            updates::download_and_install,
            updates::restart_app,
        ])
        .setup(|app| {
            sidecar::supervise(app.handle().clone());
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
