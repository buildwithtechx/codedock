import { cn } from '#/lib/utils';

export function CodedockLogo({
  className,
  size,
  withBackground = true,
}: {
  className?: string;
  size?: number;
  withBackground?: boolean;
}) {
  return (
    <svg
      viewBox="0 0 64 64"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={cn('shrink-0 select-none', className)}
      style={size ? { width: size, height: size } : undefined}
      aria-hidden="true"
    >
      {withBackground && <rect width="64" height="64" rx="16" fill="#101827" />}
      <path
        d="M18 22h17a9 9 0 0 1 9 9v2a9 9 0 0 1-9 9H29"
        stroke="#60A5FA"
        strokeWidth="6"
        strokeLinecap="round"
      />
      <path
        d="M46 42H29a9 9 0 0 1-9-9v-2a9 9 0 0 1 9-9h6"
        stroke="#A78BFA"
        strokeWidth="6"
        strokeLinecap="round"
      />
      <path d="M31 32h2" stroke="#F8FAFC" strokeWidth="6" strokeLinecap="round" />
    </svg>
  );
}
