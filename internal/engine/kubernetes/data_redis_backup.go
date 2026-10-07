package kubernetes

import (
	"bytes"
	"codedock.run/codedock/internal/models"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var redisS3Client = &http.Client{Timeout: 120 * time.Second}

func redisS3Put(ctx context.Context, dest *models.S3Destination, key string, body []byte) error {
	payload := sha256.Sum256(body)
	payloadHash := hex.EncodeToString(payload[:])
	now := time.Now().UTC()
	full := now.Format("20060102T150405Z")
	short := now.Format("20060102")
	endpoint := strings.TrimPrefix(strings.TrimPrefix(dest.Endpoint, "https://"), "http://")
	endpoint = strings.TrimRight(endpoint, "/")
	region := dest.Region
	if region == "" {
		region = "auto"
	}
	encode := func(value string) string {
		var buf strings.Builder
		for i := 0; i < len(value); i++ {
			c := value[i]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' || c == '~' {
				buf.WriteByte(c)
			} else {
				buf.WriteString(fmt.Sprintf("%%%02X", c))
			}
		}
		return buf.String()
	}
	var parts []string
	for _, p := range strings.Split(key, "/") {
		parts = append(parts, encode(p))
	}
	pathStr := "/" + encode(dest.Bucket) + "/" + strings.Join(parts, "/")
	signedHeaders := "content-type;host;x-amz-content-sha256;x-amz-date"
	canonicalHeaders := fmt.Sprintf("content-type:application/octet-stream\nhost:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n", endpoint, payloadHash, full)
	canonicalRequest := fmt.Sprintf("PUT\n%s\n\n%s\n%s\n%s", pathStr, canonicalHeaders, signedHeaders, payloadHash)
	sum := sha256.Sum256([]byte(canonicalRequest))
	scope := fmt.Sprintf("%s/%s/s3/aws4_request", short, region)
	toSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s", full, scope, hex.EncodeToString(sum[:]))
	mac := func(key []byte, data string) []byte {
		h := hmac.New(sha256.New, key)
		h.Write([]byte(data))
		return h.Sum(nil)
	}
	kDate := mac([]byte("AWS4"+dest.SecretAccessKey), short)
	kRegion := mac(kDate, region)
	kService := mac(kRegion, "s3")
	kSigning := mac(kService, "aws4_request")
	signature := hex.EncodeToString(mac(kSigning, toSign))
	urlStr := strings.TrimRight(dest.Endpoint, "/") + pathStr
	if !strings.HasPrefix(urlStr, "http") {
		urlStr = "https://" + urlStr
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, urlStr, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/octet-stream")
	request.Header.Set("x-amz-content-sha256", payloadHash)
	request.Header.Set("x-amz-date", full)
	request.Header.Set("Authorization", fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s", dest.AccessKeyID, scope, signedHeaders, signature))
	response, err := redisS3Client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		data, _ := io.ReadAll(response.Body)
		return fmt.Errorf("snapshot upload failed: %d %s", response.StatusCode, string(data))
	}
	return nil
}

func redisS3Get(ctx context.Context, dest *models.S3Destination, key string) ([]byte, error) {
	payloadHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	now := time.Now().UTC()
	full := now.Format("20060102T150405Z")
	short := now.Format("20060102")
	endpoint := strings.TrimPrefix(strings.TrimPrefix(dest.Endpoint, "https://"), "http://")
	endpoint = strings.TrimRight(endpoint, "/")
	region := dest.Region
	if region == "" {
		region = "auto"
	}
	encode := func(value string) string {
		var buf strings.Builder
		for i := 0; i < len(value); i++ {
			c := value[i]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' || c == '~' {
				buf.WriteByte(c)
			} else {
				buf.WriteString(fmt.Sprintf("%%%02X", c))
			}
		}
		return buf.String()
	}
	var parts []string
	for _, p := range strings.Split(key, "/") {
		parts = append(parts, encode(p))
	}
	pathStr := "/" + encode(dest.Bucket) + "/" + strings.Join(parts, "/")
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n", endpoint, payloadHash, full)
	canonicalRequest := fmt.Sprintf("GET\n%s\n\n%s\n%s\n%s", pathStr, canonicalHeaders, signedHeaders, payloadHash)
	sum := sha256.Sum256([]byte(canonicalRequest))
	scope := fmt.Sprintf("%s/%s/s3/aws4_request", short, region)
	toSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s", full, scope, hex.EncodeToString(sum[:]))
	mac := func(key []byte, data string) []byte {
		h := hmac.New(sha256.New, key)
		h.Write([]byte(data))
		return h.Sum(nil)
	}
	kDate := mac([]byte("AWS4"+dest.SecretAccessKey), short)
	kRegion := mac(kDate, region)
	kService := mac(kRegion, "s3")
	kSigning := mac(kService, "aws4_request")
	signature := hex.EncodeToString(mac(kSigning, toSign))
	urlStr := strings.TrimRight(dest.Endpoint, "/") + pathStr
	if !strings.HasPrefix(urlStr, "http") {
		urlStr = "https://" + urlStr
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("x-amz-content-sha256", payloadHash)
	request.Header.Set("x-amz-date", full)
	request.Header.Set("Authorization", fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s", dest.AccessKeyID, scope, signedHeaders, signature))
	response, err := redisS3Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		data, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("snapshot download failed: %d %s", response.StatusCode, string(data))
	}
	return io.ReadAll(io.LimitReader(response.Body, 512*1024*1024))
}

func (r *WorkloadRuntime) RedisSnapshot(ctx context.Context, node models.ClusterNode, record *models.ClusterData, password string, dest *models.S3Destination, key string) (int64, error) {
	namespace, name := DataIdentity(record.Spec)
	pods, err := r.dataPods(ctx, node, namespace, name)
	if err != nil || len(pods) == 0 {
		return 0, fmt.Errorf("no database pods available")
	}
	metadata, _ := pods[0]["metadata"].(map[string]any)
	pod, _ := metadata["name"].(string)
	if _, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "exec", pod, "-c", "redis", "--", "redis-cli", "-a", password, "BGSAVE"}, ""); err != nil {
		return 0, err
	}
	deadline, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	for {
		output, err := r.commands.Kubectl(deadline, node, []string{"-n", namespace, "exec", pod, "-c", "redis", "--", "redis-cli", "-a", password, "--raw", "LASTSAVE"}, "")
		if err != nil {
			return 0, err
		}
		if strings.TrimSpace(output) != "" {
			break
		}
	}
	dump, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "exec", pod, "-c", "redis", "--", "cat", "/data/dump.rdb"}, "")
	if err != nil || dump == "" {
		return 0, fmt.Errorf("snapshot capture failed")
	}
	body := []byte(dump)
	if err := redisS3Put(ctx, dest, key, body); err != nil {
		return 0, err
	}
	return int64(len(body)), nil
}

func (r *WorkloadRuntime) RedisRestore(ctx context.Context, node models.ClusterNode, record *models.ClusterData, password string, dest *models.S3Destination, key string) error {
	namespace, name := DataIdentity(record.Spec)
	body, err := redisS3Get(ctx, dest, key)
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return fmt.Errorf("snapshot archive is empty")
	}
	pods, err := r.dataPods(ctx, node, namespace, name)
	if err != nil || len(pods) == 0 {
		return fmt.Errorf("no database pods available")
	}
	for _, pod := range pods {
		metadata, _ := pod["metadata"].(map[string]any)
		podName, _ := metadata["name"].(string)
		_, _ = r.commands.Kubectl(ctx, node, []string{"-n", namespace, "exec", podName, "-c", "redis", "--", "redis-cli", "-a", password, "SHUTDOWN", "NOSAVE"}, "")
		if _, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "exec", "-i", podName, "-c", "redis", "--", "sh", "-c", "cat > /data/dump.rdb"}, string(body)); err != nil {
			return err
		}
		if _, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "delete", "pod", podName, "--wait=true", "--timeout=180s"}, ""); err != nil {
			return err
		}
	}
	return r.DataObserved(ctx, node, record)
}

func (r *WorkloadRuntime) RedisFailover(ctx context.Context, node models.ClusterNode, record *models.ClusterData) error {
	namespace, name := DataIdentity(record.Spec)
	if _, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "delete", "pod", name + "-0", "--wait=true", "--timeout=180s"}, ""); err != nil {
		return err
	}
	return r.DataObserved(ctx, node, record)
}
