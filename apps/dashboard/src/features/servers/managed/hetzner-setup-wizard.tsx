import {
  HETZNER_WIZARD_STEPS as steps,
  type HetznerWizardStep as WizardStep,
} from '@codedock/setup';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import type { ManagedCredential, ManagedOffers, ManagedQuota } from '#/interfaces/managed';
import { cn } from '#/lib/utils';
import { managedService } from '#/services/managed';
import { useOrganizationStore } from '#/stores/organization-store';
import { HetznerCredentialStep } from './hetzner-credential-step';
import { HetznerPlanStep } from './hetzner-plan-step';
import { HetznerProgressView } from './hetzner-progress-view';
import { HetznerReviewStep } from './hetzner-review-step';

export function HetznerSetupWizard({ onDone }: { onDone: (serverId: string) => void }) {
  const organizationId = useOrganizationStore((state) => state.activeOrganizationId);
  const [step, setStep] = useState<WizardStep>('credential');
  const [credentials, setCredentials] = useState<ManagedCredential[]>([]);
  const [credentialsLoading, setCredentialsLoading] = useState(true);
  const [selectedCredentialId, setSelectedCredentialId] = useState<string | null>(null);
  const [creatingCredential, setCreatingCredential] = useState(false);
  const [createError, setCreateError] = useState<string | null>(null);

  const [offers, setOffers] = useState<ManagedOffers | null>(null);
  const [offersLoading, setOffersLoading] = useState(false);
  const [offersError, setOffersError] = useState<string | null>(null);
  const [location, setLocation] = useState('');
  const [tierName, setTierName] = useState('');

  const [name, setName] = useState('');
  const [image, setImage] = useState('');
  const [quota, setQuota] = useState<ManagedQuota | null>(null);
  const [quotaLoading, setQuotaLoading] = useState(false);

  const [working, setWorking] = useState(false);
  const [reviewError, setReviewError] = useState<string | null>(null);
  const [operationId, setOperationId] = useState<string | null>(null);
  const [serverId, setServerId] = useState<string | null>(null);

  useEffect(() => {
    if (!organizationId) {
      setCredentialsLoading(false);
      return;
    }
    let cancelled = false;
    setCredentialsLoading(true);
    setQuotaLoading(true);
    managedService
      .listCredentials(organizationId)
      .then((rows) => {
        if (cancelled) return;
        setCredentials(rows);
        if (rows.length === 1) setSelectedCredentialId(rows[0].id);
      })
      .catch(() => {
        if (!cancelled) setCredentials([]);
      })
      .finally(() => {
        if (!cancelled) setCredentialsLoading(false);
      });
    managedService
      .getQuota(organizationId)
      .then((value) => {
        if (!cancelled) setQuota(value);
      })
      .catch(() => {
        if (!cancelled) setQuota(null);
      })
      .finally(() => {
        if (!cancelled) setQuotaLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [organizationId]);

  useEffect(() => {
    if (!organizationId || !selectedCredentialId) {
      setOffers(null);
      return;
    }
    let cancelled = false;
    setOffersLoading(true);
    setOffersError(null);
    managedService
      .offers(organizationId, selectedCredentialId)
      .then((value) => {
        if (cancelled) return;
        setOffers(value);
        setLocation((prev) =>
          value.locations.some((entry) => entry.name === prev)
            ? prev
            : (value.locations[0]?.name ?? '')
        );
        setTierName((prev) =>
          value.offers.some((offer) => offer.tier.name === prev)
            ? prev
            : (value.offers[0]?.tier.name ?? '')
        );
        setImage((prev) => (value.images.includes(prev) ? prev : (value.images[0] ?? '')));
      })
      .catch((err: unknown) => {
        if (cancelled) return;
        setOffers(null);
        setOffersError(err instanceof Error ? err.message : 'Failed to load Hetzner offers');
      })
      .finally(() => {
        if (!cancelled) setOffersLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [organizationId, selectedCredentialId]);

  if (!organizationId) {
    return (
      <div className="rounded-xl border border-border/60 p-6 text-muted-foreground text-sm">
        Select an organization first — managed servers are billed and quota-checked per
        organization.
      </div>
    );
  }

  const selectedOffer = offers?.offers.find((offer) => offer.tier.name === tierName) ?? null;
  const stepIndex = steps.findIndex((entry) => entry.id === step);

  const handleCreateCredential = async (label: string, token: string) => {
    setCreatingCredential(true);
    setCreateError(null);
    try {
      const credential = await managedService.createCredential(organizationId, {
        provider: 'hetzner',
        label,
        token,
      });
      setCredentials((prev) => [...prev, credential]);
      setSelectedCredentialId(credential.id);
      toast.success('Hetzner token verified and saved');
    } catch (err) {
      setCreateError(err instanceof Error ? err.message : 'Token verification failed');
    } finally {
      setCreatingCredential(false);
    }
  };

  const handleConfirm = async () => {
    if (!selectedCredentialId || !selectedOffer) return;
    setWorking(true);
    setReviewError(null);
    try {
      const review = await managedService.reviewProvision(organizationId, {
        credentialId: selectedCredentialId,
        name: name.trim(),
        region: location,
        serverType: selectedOffer.tier.name,
        image,
      });
      await managedService.applyOperation(review.review.operation.id, review.review.confirmation);
      setOperationId(review.review.operation.id);
      setServerId(review.plan.serverId);
      setStep('progress');
    } catch (err) {
      setReviewError(err instanceof Error ? err.message : 'Provision review failed');
    } finally {
      setWorking(false);
    }
  };

  return (
    <div className="space-y-6">
      <ol className="flex items-center gap-2">
        {steps.map((entry, index) => (
          <li key={entry.id} className="flex flex-1 items-center gap-2 last:flex-none">
            <span
              className={cn(
                'flex size-6 shrink-0 items-center justify-center rounded-full font-medium text-[11px]',
                index < stepIndex && 'bg-primary text-primary-foreground',
                index === stepIndex && 'bg-primary text-primary-foreground',
                index > stepIndex && 'bg-muted text-muted-foreground'
              )}
            >
              {index + 1}
            </span>
            <span
              className={cn(
                'hidden text-xs sm:block',
                index <= stepIndex ? 'font-medium text-foreground' : 'text-muted-foreground'
              )}
            >
              {entry.label}
            </span>
            {index < steps.length - 1 && <span className="mx-1 h-px flex-1 bg-border/70" />}
          </li>
        ))}
      </ol>

      {step === 'credential' && (
        <div className="space-y-4">
          <HetznerCredentialStep
            credentials={credentials}
            selectedId={selectedCredentialId}
            loading={credentialsLoading}
            creating={creatingCredential}
            createError={createError}
            onSelect={setSelectedCredentialId}
            onCreate={(label, token) => void handleCreateCredential(label, token)}
          />
          <div className="flex items-center justify-end border-border/60 border-t pt-4">
            <Button
              type="button"
              onClick={() => setStep('plan')}
              disabled={!selectedCredentialId}
              className="text-xs"
            >
              Continue
            </Button>
          </div>
        </div>
      )}

      {step === 'plan' && (
        <div className="space-y-4">
          <HetznerPlanStep
            offers={offers}
            loading={offersLoading}
            loadError={offersError}
            location={location}
            tierName={tierName}
            onLocationChange={setLocation}
            onTierChange={setTierName}
          />
          <div className="flex items-center justify-between border-border/60 border-t pt-4">
            <Button
              type="button"
              variant="ghost"
              onClick={() => setStep('credential')}
              className="text-xs"
            >
              Back
            </Button>
            <Button
              type="button"
              onClick={() => setStep('review')}
              disabled={!location || !selectedOffer}
              className="text-xs"
            >
              Continue
            </Button>
          </div>
        </div>
      )}

      {step === 'review' && selectedOffer && (
        <HetznerReviewStep
          name={name}
          image={image}
          images={offers?.images ?? []}
          location={location}
          tierName={selectedOffer.tier.name}
          tierLabel={`${selectedOffer.tier.cores} vCPU · ${selectedOffer.tier.memoryGB} GB · ${selectedOffer.tier.diskGB} GB`}
          monthlyPrice={selectedOffer.prices[location] ?? selectedOffer.tier.monthlyPrice}
          quota={quota}
          quotaLoading={quotaLoading}
          working={working}
          error={reviewError}
          onNameChange={setName}
          onImageChange={setImage}
          onBack={() => setStep('plan')}
          onConfirm={() => void handleConfirm()}
        />
      )}

      {step === 'progress' && operationId && serverId && (
        <HetznerProgressView
          operationId={operationId}
          serverName={name.trim()}
          serverId={serverId}
          onDone={onDone}
        />
      )}
    </div>
  );
}
