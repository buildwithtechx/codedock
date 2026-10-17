import { Check, Loader2, Lock, Mail, User } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from '@codedock/ui/dialog';
import { Input } from '@codedock/ui/input';
import { InputOTP, InputOTPGroup, InputOTPSlot } from '@codedock/ui/input-otp';
import {
  useChangePassword,
  useGetProfile,
  useRequestEmailChange,
  useUpdateProfile,
  useVerifyEmailChange,
} from '#/features/profile';
import { SettingsRow } from '#/features/settings/settings-row';
import { SettingsSection } from '#/features/settings/settings-section';

export function ProfileNameForm() {
  const { data: profile, isLoading } = useGetProfile();
  const updateProfile = useUpdateProfile();
  const [name, setName] = useState('');

  if (!isLoading && profile?.data && !name && profile.data.name) {
    setName(profile.data.name);
  }

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    updateProfile.mutate(
      { name },
      {
        onSuccess: () => toast.success('Profile name updated!'),
        onError: (err) => toast.error(err.message),
      }
    );
  };

  return (
    <form onSubmit={handleSave}>
      <SettingsSection
        icon={<User className="size-4 text-primary" />}
        title="Profile Name"
        description="Your public display name on Codedock."
        action={
          <Button
            type="submit"
            size="sm"
            disabled={isLoading || updateProfile.isPending || name === profile?.data?.name}
          >
            {updateProfile.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Check className="mr-2 h-4 w-4" />
            )}
            Save Changes
          </Button>
        }
      >
        <SettingsRow
          label="Full Name"
          description="Update your display name across workspace views."
          bordered={false}
        >
          <Input
            id="name"
            placeholder="John Doe"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full text-sm"
          />
        </SettingsRow>
      </SettingsSection>
    </form>
  );
}

export function ProfileEmailForm() {
  const { data: profile, isLoading } = useGetProfile();
  const requestEmailChange = useRequestEmailChange();
  const verifyEmailChange = useVerifyEmailChange();

  const [email, setEmail] = useState('');
  const [otpOpen, setOtpOpen] = useState(false);
  const [otp, setOtp] = useState('');

  if (!isLoading && profile?.data && !email && profile.data.email) {
    setEmail(profile.data.email);
  }

  const handleRequest = (e: React.FormEvent) => {
    e.preventDefault();
    requestEmailChange.mutate(
      { newEmail: email },
      {
        onSuccess: () => {
          setOtpOpen(true);
          toast.success('Verification code sent to your new email.');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  const handleVerify = (e: React.FormEvent) => {
    e.preventDefault();
    verifyEmailChange.mutate(
      { otp },
      {
        onSuccess: () => {
          setOtpOpen(false);
          setOtp('');
          toast.success('Email updated successfully!');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  return (
    <>
      <form onSubmit={handleRequest}>
        <SettingsSection
          icon={<Mail className="size-4 text-primary" />}
          title="Email Address"
          description="Manage the email address associated with your account."
          action={
            <Button
              type="submit"
              size="sm"
              disabled={isLoading || requestEmailChange.isPending || email === profile?.data?.email}
            >
              {requestEmailChange.isPending ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <Check className="mr-2 h-4 w-4" />
              )}
              Save Changes
            </Button>
          }
        >
          <SettingsRow
            label="New Email Address"
            description="Change the email address associated with your account."
            bordered={false}
          >
            <Input
              id="email"
              type="email"
              placeholder="john@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full text-sm"
            />
          </SettingsRow>
        </SettingsSection>
      </form>

      <Dialog open={otpOpen} onOpenChange={setOtpOpen}>
        <DialogContent className="gap-0 border-border/50 bg-card/95 p-0 backdrop-blur-xl sm:max-w-100 [&>button]:hidden">
          <div className="px-5 pt-5 pb-4">
            <div className="flex items-start justify-between">
              <div className="flex flex-col">
                <DialogTitle className="flex items-center gap-2 font-bold text-foreground text-xl tracking-tight">
                  <Mail className="h-5 w-5 text-primary" />
                  Verify Email Change
                </DialogTitle>
                <DialogDescription>Code sent to new email</DialogDescription>
              </div>
              <DialogClose asChild>
                <Button
                  variant="ghost"
                  className="font-medium text-foreground/80 text-sm hover:bg-transparent hover:text-foreground"
                >
                  CLOSE
                </Button>
              </DialogClose>
            </div>
          </div>
          <div className="h-px w-full bg-border/50" />
          <div className="flex flex-col items-center gap-4 p-5">
            <InputOTP maxLength={6} value={otp} onChange={setOtp}>
              <InputOTPGroup>
                <InputOTPSlot index={0} />
                <InputOTPSlot index={1} />
                <InputOTPSlot index={2} />
                <InputOTPSlot index={3} />
                <InputOTPSlot index={4} />
                <InputOTPSlot index={5} />
              </InputOTPGroup>
            </InputOTP>
          </div>
          <div className="flex items-center justify-end gap-3 p-5 pt-0">
            <Button
              variant="ghost"
              onClick={() => setOtpOpen(false)}
              className="h-9 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              Cancel
            </Button>
            <Button
              onClick={handleVerify}
              disabled={otp.length !== 6 || verifyEmailChange.isPending}
              className="h-9 gap-2 font-mono font-semibold text-[11px] uppercase tracking-wider"
            >
              {verifyEmailChange.isPending ? (
                <Loader2 className="h-3.5 w-3.5 animate-spin" />
              ) : (
                <Check className="h-3.5 w-3.5" />
              )}
              Verify & Update
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}

export function ProfilePasswordForm() {
  const changePassword = useChangePassword();
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');

  const handleSave = (e: React.FormEvent) => {
    e.preventDefault();
    changePassword.mutate(
      { oldPassword, newPassword },
      {
        onSuccess: () => {
          toast.success('Password updated successfully!');
          setOldPassword('');
          setNewPassword('');
        },
        onError: (err) => toast.error(err.message),
      }
    );
  };

  return (
    <form onSubmit={handleSave}>
      <SettingsSection
        icon={<Lock className="size-4 text-primary" />}
        title="Change Password"
        description="Ensure your account is using a secure password."
        action={
          <Button
            type="submit"
            size="sm"
            disabled={changePassword.isPending || !oldPassword || !newPassword}
          >
            {changePassword.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Check className="mr-2 h-4 w-4" />
            )}
            Save Changes
          </Button>
        }
      >
        <SettingsRow label="Current Password" description="Update the password you use to sign in.">
          <Input
            id="oldPassword"
            type="password"
            value={oldPassword}
            onChange={(e) => setOldPassword(e.target.value)}
            className="w-full text-sm"
            placeholder="Current Password"
          />
        </SettingsRow>
        <SettingsRow
          label="New Password"
          description="Must be at least 8 characters long."
          bordered={false}
        >
          <Input
            id="newPassword"
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            className="w-full text-sm"
            placeholder="New Password"
          />
        </SettingsRow>
      </SettingsSection>
    </form>
  );
}
