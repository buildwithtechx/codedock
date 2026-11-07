export type MailProviderId =
  | 'custom'
  | 'resend'
  | 'ses'
  | 'oracle'
  | 'sendgrid'
  | 'mailgun'
  | 'postmark'
  | 'gmail'
  | 'fastmail';

export interface MailProvider {
  id: MailProviderId;
  label: string;
  logoSrc?: string;
  imapHost: string;
  imapPort: number;
  smtpHost: string;
  smtpPort: number;
  sendOnly?: boolean;
  hint?: string;
}

const brandLogo = (domain: string) => `https://www.google.com/s2/favicons?domain=${domain}&sz=128`;

export const MAIL_PROVIDERS: readonly MailProvider[] = [
  {
    id: 'custom',
    label: 'Custom',
    imapHost: '',
    imapPort: 993,
    smtpHost: '',
    smtpPort: 465,
  },
  {
    id: 'gmail',
    label: 'Gmail',
    logoSrc: 'https://ssl.gstatic.com/images/branding/product/2x/gmail_2020q4_64dp.png',
    imapHost: 'imap.gmail.com',
    imapPort: 993,
    smtpHost: 'smtp.gmail.com',
    smtpPort: 465,
    hint: 'Use an app password (Google Account → Security → App passwords) with 2FA on.',
  },
  {
    id: 'ses',
    label: 'Amazon SES',
    logoSrc: brandLogo('aws.amazon.com'),
    sendOnly: true,
    imapHost: '',
    imapPort: 993,
    smtpHost: 'email-smtp.us-east-1.amazonaws.com',
    smtpPort: 587,
    hint: 'Use the SMTP credentials from the SES console (an IAM SMTP user).',
  },
  {
    id: 'oracle',
    label: 'Oracle Cloud',
    logoSrc: brandLogo('oracle.com'),
    sendOnly: true,
    imapHost: '',
    imapPort: 993,
    smtpHost: 'smtp.email.us-ashburn-1.oci.oraclecloud.com',
    smtpPort: 587,
    hint: 'OCI Email Delivery: generate SMTP credentials in the OCI console.',
  },
  {
    id: 'resend',
    label: 'Resend',
    logoSrc: brandLogo('resend.com'),
    sendOnly: true,
    imapHost: '',
    imapPort: 993,
    smtpHost: 'smtp.resend.com',
    smtpPort: 465,
    hint: 'Username is "resend"; paste your Resend API key as the password.',
  },
  {
    id: 'sendgrid',
    label: 'SendGrid',
    logoSrc: brandLogo('sendgrid.com'),
    sendOnly: true,
    imapHost: '',
    imapPort: 993,
    smtpHost: 'smtp.sendgrid.net',
    smtpPort: 587,
    hint: 'Username is "apikey"; paste your SendGrid API key as the password.',
  },
  {
    id: 'mailgun',
    label: 'Mailgun',
    logoSrc: brandLogo('mailgun.com'),
    sendOnly: true,
    imapHost: '',
    imapPort: 993,
    smtpHost: 'smtp.mailgun.org',
    smtpPort: 587,
    hint: 'Use your Mailgun SMTP credentials (postmaster@your-domain + its SMTP password).',
  },
  {
    id: 'postmark',
    label: 'Postmark',
    logoSrc: brandLogo('postmarkapp.com'),
    sendOnly: true,
    imapHost: '',
    imapPort: 993,
    smtpHost: 'smtp.postmarkapp.com',
    smtpPort: 587,
    hint: 'Use your Postmark Server API token as both username and password.',
  },
  {
    id: 'fastmail',
    label: 'Fastmail',
    logoSrc: brandLogo('fastmail.com'),
    imapHost: 'imap.fastmail.com',
    imapPort: 993,
    smtpHost: 'smtp.fastmail.com',
    smtpPort: 465,
    hint: 'Create an app password in Fastmail → Settings → Privacy & Security.',
  },
];

export function getMailProvider(id: MailProviderId): MailProvider {
  return MAIL_PROVIDERS.find((p) => p.id === id) ?? MAIL_PROVIDERS[0];
}
