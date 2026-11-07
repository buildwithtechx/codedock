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
import { Check, Loader2, Lock, Mail, User } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { SettingsRow } from '#/features/settings/settings-row';
import { SettingsSection } from '#/features/settings/settings-section';
import {
  useChangePassword,
  useGetProfile,
  useRequestEmailChange,
  useUpdateProfile,
  useVerifyEmailChange,
} from '#/hooks/use-profile';

export function ProfileDetailsForm() {
  const { data: profile, isLoading } = useGetProfile();
  const updateProfile = useUpdateProfile();
  const requestEmailChange = useRequestEmailChange();
  const verifyEmailChange = useVerifyEmailChange();

  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [otpOpen, setOtpOpen] = useState(false);
  const [otp, setOtp] = useState('');

  useEffect(() => {
    if (profile?.data) {
      if (profile.data.name !== undefined) setName(profile.data.name);
      if (profile.data.email !== undefined) setEmail(profile.data.email);
    }
  }, [profile?.data]);

  const originalName = profile?.data?.name ?? '';
  const originalEmail = profile?.data?.email ?? '';
  const isNameChanged = name.trim() !== originalName;
  const isEmailChanged = email.trim().toLowerCase() !== originalEmail.toLowerCase();
  const isPending = updateProfile.isPending || requestEmailChange.isPending;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Profile name cannot be empty');
      return;
    }
    if (!email.trim()) {
      toast.error('Email address cannot be empty');
      return;
    }

    if (isNameChanged && !isEmailChanged) {
      updateProfile.mutate(
        { name: name.trim() },
        {
          onSuccess: () => toast.success('Profile updated successfully'),
          onError: (err) => toast.error(err.message),
        }
      );
    } else if (isEmailChanged) {
      if (isNameChanged) {
        updateProfile.mutate({ name: name.trim() });
      }
      requestEmailChange.mutate(
        { newEmail: email.trim() },
        {
          onSuccess: () => {
            setOtpOpen(true);
            toast.success('Verification code sent to your new email');
          },
          onError: (err) => toast.error(err.message),
        }
      );
    }
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
      <form onSubmit={handleSubmit}>
        <SettingsSection
          icon={<User className="size-4" />}
          title="Profile Information"
          description="Update your personal account details and primary email address."
        >
          <SettingsRow
            label="Full Name"
            description="Shown on your profile and across team workspaces."
          >
            <Input
              id="profile-name"
              placeholder="John Doe"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full text-sm"
            />
          </SettingsRow>

          <SettingsRow
            label="Email Address"
            description="The primary email address used for login and notifications."
            bordered={false}
          >
            <Input
              id="profile-email"
              type="email"
              placeholder="john@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full text-sm"
            />
          </SettingsRow>

          <div className="flex justify-end border-border/40 border-t pt-4">
            <Button
              type="submit"
              size="sm"
              disabled={isLoading || isPending || (!isNameChanged && !isEmailChanged)}
              className="gap-1.5"
            >
              {isPending ? (
                <Loader2 className="h-3.5 w-3.5 animate-spin" />
              ) : (
                <Check className="h-3.5 w-3.5" />
              )}
              <span>Save Profile</span>
            </Button>
          </div>
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
        icon={<Lock className="size-4" />}
        title="Password & Security"
        description="Ensure your account is using a secure and robust password."
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

        <div className="flex justify-end border-border/40 border-t pt-4">
          <Button
            type="submit"
            size="sm"
            disabled={changePassword.isPending || !oldPassword || !newPassword}
            className="gap-1.5"
          >
            {changePassword.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Check className="mr-2 h-4 w-4" />
            )}
            <span>Update Password</span>
          </Button>
        </div>
      </SettingsSection>
    </form>
  );
}

export const ProfileNameForm = ProfileDetailsForm;
export const ProfileEmailForm = ProfileDetailsForm;
