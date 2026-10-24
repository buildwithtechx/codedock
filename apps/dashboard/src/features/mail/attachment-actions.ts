import { toast } from 'sonner';
import { escapeHtml, getOfficePreviewKind, renderOfficePreview } from './office-preview';

export interface MailAttachment {
  body: string;
  mimeType: string;
  filename: string;
  attachmentId: string;
}

const INLINE_SAFE_TYPES = new Set([
  'application/pdf',
  'image/png',
  'image/jpeg',
  'image/gif',
  'image/webp',
  'image/avif',
  'text/plain',
]);

export const getInlineAttachmentType = (mimeType: string): string | null => {
  const normalized = mimeType.toLowerCase().split(';')[0]?.trim() ?? '';
  return INLINE_SAFE_TYPES.has(normalized) ? normalized : null;
};

export const downloadAttachment = async (attachment: MailAttachment): Promise<void> => {
  try {
    const byteCharacters = atob(attachment.body);
    const byteArray = new Uint8Array(byteCharacters.length);
    for (let i = 0; i < byteCharacters.length; i++) {
      byteArray[i] = byteCharacters.charCodeAt(i);
    }
    const blob = new Blob([byteArray], { type: attachment.mimeType });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = attachment.filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error('Error downloading attachment:', error);
    toast.error('Failed to download attachment');
  }
};

export const openAttachmentViewer = (): Window | null => {
  const width = 800;
  const height = 600;
  const left = (window.screen.width - width) / 2;
  const top = (window.screen.height - height) / 2;
  return window.open(
    '',
    'attachment-viewer',
    `width=${width},height=${height},left=${left},top=${top},resizable=yes,scrollbars=yes,status=no,location=no,menubar=no`
  );
};

const sandboxedViewer = (title: string, html: string): string =>
  `<!doctype html><html><head><meta charset="utf-8"><meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src 'unsafe-inline'"><title>${escapeHtml(
    title
  )}</title><style>html,body,iframe{display:block;width:100%;height:100%;margin:0;border:0}</style></head><body><iframe sandbox srcdoc="${escapeHtml(
    html
  )}"></iframe></body></html>`;

const viewerBlob = async (
  attachment: { mimeType: string; filename: string },
  bytes: Uint8Array
): Promise<Blob | null> => {
  const officeKind = getOfficePreviewKind(attachment.mimeType, attachment.filename);
  if (officeKind) {
    const preview = await renderOfficePreview(officeKind, bytes, attachment.filename).catch(
      () => null
    );
    return preview
      ? new Blob([sandboxedViewer(attachment.filename, preview)], { type: 'text/html' })
      : null;
  }
  const inlineType = getInlineAttachmentType(attachment.mimeType);
  return inlineType ? new Blob([bytes as BlobPart], { type: inlineType }) : null;
};

export const openAttachment = async (
  attachment: MailAttachment,
  pendingPopup?: Window | null
): Promise<void> => {
  try {
    const byteCharacters = atob(attachment.body);
    const byteArray = new Uint8Array(byteCharacters.length);
    for (let i = 0; i < byteCharacters.length; i++) {
      byteArray[i] = byteCharacters.charCodeAt(i);
    }

    const blob = await viewerBlob(attachment, byteArray);
    if (!blob) {
      pendingPopup?.close();
      await downloadAttachment(attachment);
      return;
    }

    const url = window.URL.createObjectURL(blob);
    const popup = pendingPopup === undefined ? openAttachmentViewer() : pendingPopup;

    if (popup) {
      popup.location.href = url;
      popup.focus();
      setTimeout(() => window.URL.revokeObjectURL(url), 1000);
    } else {
      window.URL.revokeObjectURL(url);
    }
  } catch (error) {
    pendingPopup?.close();
    console.error('Error opening attachment:', error);
    toast.error('Failed to open attachment');
  }
};
