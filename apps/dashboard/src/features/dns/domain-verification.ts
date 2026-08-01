import { useState } from 'react';
import { toast } from 'sonner';
import { useVerifyDomain } from './hooks';
import type { DomainVerifyResult } from './interfaces';

interface DomainSummary {
  id: string;
}

function toastForVerificationResult(result: DomainVerifyResult) {
  if (result.status === 'resolves_to_server') {
    toast.success(result.message);
    return;
  }

  if (result.status === 'resolves_to_different_ip') {
    toast.warning(result.message);
    return;
  }

  toast.error(result.message);
}

export function useDomainVerification() {
  const verifyDomain = useVerifyDomain();
  const [verifyingMap, setVerifyingMap] = useState<Record<string, boolean>>({});
  const [verificationResults, setVerificationResults] = useState<
    Record<string, DomainVerifyResult>
  >({});
  const [isVerifyingAll, setIsVerifyingAll] = useState(false);

  const runVerification = async (domainId: string, notify: boolean) => {
    const res = await verifyDomain.mutateAsync(domainId);
    const data = res.data;

    if (!data) {
      if (notify) {
        toast.error('Verification failed');
      }
      return undefined;
    }

    setVerificationResults((prev) => ({ ...prev, [domainId]: data }));

    if (notify) {
      toastForVerificationResult(data);
    }

    return data;
  };

  const verifyOne = async (domainId: string) => {
    setVerifyingMap((prev) => ({ ...prev, [domainId]: true }));

    try {
      return await runVerification(domainId, true);
    } catch {
      toast.error('Verification failed');
      return undefined;
    } finally {
      setVerifyingMap((prev) => ({ ...prev, [domainId]: false }));
    }
  };

  const verifyAll = async (domains: DomainSummary[], concurrency = 4) => {
    if (domains.length === 0) {
      return { verifiedCount: 0, failedCount: 0, totalCount: 0 };
    }

    const workerCount = Math.max(1, Math.min(concurrency, domains.length));
    setIsVerifyingAll(true);
    let verifiedCount = 0;
    let failedCount = 0;
    let index = 0;

    const worker = async () => {
      while (index < domains.length) {
        const currentIndex = index;
        index += 1;
        const domain = domains[currentIndex];

        setVerifyingMap((prev) => ({ ...prev, [domain.id]: true }));

        try {
          const data = await runVerification(domain.id, false);
          if (data?.verified) {
            verifiedCount += 1;
          } else {
            failedCount += 1;
          }
        } catch {
          failedCount += 1;
        } finally {
          setVerifyingMap((prev) => ({ ...prev, [domain.id]: false }));
        }
      }
    };

    try {
      await Promise.all(Array.from({ length: workerCount }, () => worker()));
      return {
        verifiedCount,
        failedCount,
        totalCount: domains.length,
      };
    } finally {
      setIsVerifyingAll(false);
    }
  };

  return {
    isVerifyingAll,
    verificationResults,
    verifyingMap,
    verifyAll,
    verifyOne,
  };
}
