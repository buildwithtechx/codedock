export type SupportCategory =
  | 'deployment'
  | 'billing'
  | 'account'
  | 'general'
  | 'technical'
  | 'urgent';

export type SupportStatus = 'open' | 'closed' | 'resolved' | 'pending_user';

export interface SupportMessage {
  id: string;
  sender: 'user' | 'staff';
  senderName: string;
  content: string;
  createdAt: string;
}

export interface SupportTicket {
  id: string;
  subject: string;
  category: SupportCategory;
  status: SupportStatus;
  hasDiagnosticBundle?: boolean;
  messages: SupportMessage[];
  createdAt: string;
  updatedAt: string;
}

export interface CreateTicketInput {
  subject: string;
  category: SupportCategory;
  message: string;
  attachDiagnostics?: boolean;
}

export interface CloudSupportSession {
  account?: {
    id: string;
    key: string;
    email: string;
    name?: string | null;
  } | null;
}
