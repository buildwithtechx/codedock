# clickhouse is intentionally absent: its bundled tile is hand-crafted for
# cross-theme contrast (the bare brand mark is near-invisible on light cards).
$pairs = @(
  'calcom:caldotcom',
  'directus:directus','freshrss:freshrss',
  'ghost:ghost','gitea:gitea','grafana:grafana',
  'kafka:apachekafka','mariadb:mariadb',
  'meilisearch:meilisearch','metabase:metabase','minio:minio',
  'mongodb:mongodb','mysql:mysql','n8n:n8n','nats:natsdotio',
  'nextcloud:nextcloud','pocketbase:pocketbase','postgres:postgresql',
  'qdrant:qdrant','rabbitmq:rabbitmq','redis:redis',
  'supabase-studio:supabase','timescaledb:timescale','umami:umami',
  'uptimekuma:uptimekuma','vaultwarden:vaultwarden','wordpress:wordpress'
)
$ok = @()
$miss = @()
foreach ($p in $pairs) {
  $k, $s = $p.Split(':')
  $dest = "apps/dashboard/public/app-logos/$k.svg"
  try {
    Invoke-WebRequest -Uri "https://cdn.simpleicons.org/$s" -OutFile $dest -TimeoutSec 20
    $content = Get-Content -Raw $dest
    if ($content.StartsWith('<svg')) { $ok += $k } else { Remove-Item $dest; $miss += $p }
  } catch {
    $miss += $p
  }
}
Write-Output ("ok=" + $ok.Count + " miss=" + $miss.Count)
Write-Output ($miss -join ', ')
