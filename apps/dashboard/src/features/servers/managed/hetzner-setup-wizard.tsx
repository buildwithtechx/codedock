import {
  HETZNER_WIZARD_STEPS as steps,
  type HetznerWizardStep as WizardStep,
  nextStep,
  prevStep,
} from '@codedock/setup';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { cn } from '#/lib/utils';
import {
  useApplyManagedOperation,
  useCreateManagedCredential,
  useManagedCredentials,
  useManagedOffers,
  useManagedQuota,
  useReviewManagedProvision,
} from '#/hooks/use-managed';
import { useOrganizationStore } from '#/stores/organization-store';
import { HetznerCredentialStep } from './hetzner-credential-step';
import { HetznerPlanStep } from './hetzner-plan-step';
import { HetznerProgressView } from './hetzner-progress-view';
import { HetznerReviewStep } from './hetzner-review-step';

export function HetznerSetupWizard({ onDone }: { onDone: (serverId: string) => void }) {
  const organizationId = useOrganizationStore((state) => state.activeOrganizationId);
  const [step, setStep] = useState<WizardStep>('credential');
  const [selectedCredentialId, setSelectedCredentialId] = useState<string | null>(null);
  const [createError, setCreateError] = useState<string | null>(null);

  const { data: credentialsData = [], isLoading: credentialsLoading } = useManagedCredentials(
    organizationId ?? undefined
  );
  const { data: quotaData = null, isLoading: quotaLoading } = useManagedQuota(
    organizationId ?? undefined
  );
  const {
    data: offersData = null,
    isLoading: offersLoading,
    error: offersQueryError,
  } = useManagedOffers(organizationId ?? undefined, selectedCredentialId ?? undefined);

  const createCredentialMutation = useCreateManagedCredential();
  const reviewProvisionMutation = useReviewManagedProvision();
  const applyOperationMutation = useApplyManagedOperation();

  const [location, setLocation] = useState('');
  const [tierName, setTierName] = useState('');
  const [name, setName] = useState('');
  const [image, setImage] = useState('');

  const [working, setWorking] = useState(false);
  const [reviewError, setReviewError] = useState<string | null>(null);
  const [operationId, setOperationId] = useState<string | null>(null);
  const [serverId, setServerId] = useState<string | null>(null);

  useEffect(() => {
    if (credentialsData.length === 1 && !selectedCredentialId) {
      setSelectedCredentialId(credentialsData[0].id);
    }
  }, [credentialsData, selectedCredentialId]);

  useEffect(() => {
    if (offersData) {
      setLocation((prev) =>
        offersData.locations.some((entry) => entry.name === prev)
          ? prev
          : (offersData.locations[0]?.name ?? '')
      );
      setTierName((prev) =>
        offersData.offers.some((offer) => offer.tier.name === prev)
          ? prev
          : (offersData.offers[0]?.tier.name ?? '')
      );
      setImage((prev) => (offersData.images.includes(prev) ? prev : (offersData.images[0] ?? '')));
    }
  }, [offersData]);

  if (!organizationId) {
    return (
      <div className="rounded-xl border border-border/60 p-6 text-muted-foreground text-sm">
        Select an organization first — managed servers are billed and quota-checked per
        organization.
      </div>
    );
  }

  const offers = offersData;
  const offersError = offersQueryError instanceof Error ? offersQueryError.message : null;
  const credentials = credentialsData;
  const quota = quotaData;
  const creatingCredential = createCredentialMutation.isPending;

  const selectedOffer = offers?.offers.find((offer) => offer.tier.name === tierName) ?? null;
  const stepIndex = steps.findIndex((entry) => entry.id === step);
  const stepIds = steps.map((entry) => entry.id);

  const handleNext = () => {
    const next = nextStep(stepIds, step);
    if (next) setStep(next);
  };

  const handleBack = () => {
    const prev = prevStep(stepIds, step);
    if (prev) setStep(prev);
  };

  const handleCreateCredential = async (label: string, token: string) => {
    if (!organizationId) return;
    setCreateError(null);
    try {
      const credential = await createCredentialMutation.mutateAsync({
        organizationId,
        payload: { provider: 'hetzner', label, token },
      });
      setSelectedCredentialId(credential.id);
      toast.success('Hetzner token verified and saved');
    } catch (err) {
      setCreateError(err instanceof Error ? err.message : 'Token verification failed');
    }
  };

  const handleConfirm = async () => {
    if (!organizationId || !selectedCredentialId || !selectedOffer) return;
    setWorking(true);
    setReviewError(null);
    try {
      const review = await reviewProvisionMutation.mutateAsync({
        organizationId,
        payload: {
          credentialId: selectedCredentialId,
          name: name.trim(),
          region: location,
          serverType: selectedOffer.tier.name,
          image,
        },
      });
      await applyOperationMutation.mutateAsync({
        operationId: review.review.operation.id,
        confirmation: review.review.confirmation,
      });
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
              onClick={handleNext}
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
            <Button type="button" variant="ghost" onClick={handleBack} className="text-xs">
              Back
            </Button>
            <Button
              type="button"
              onClick={handleNext}
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
          onBack={handleBack}
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
