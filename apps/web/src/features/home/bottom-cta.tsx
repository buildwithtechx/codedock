import { InstallCommand } from '../../components/install-command';
import { productLinks } from '../../lib/product-links';
export function BottomCta() {
  return (
    <section className="relative overflow-hidden px-6 py-24 text-center">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_at_center,var(--color-primary),transparent_65%)] opacity-10" />
      <div className="relative mx-auto max-w-3xl">
        <h2 className="font-bold text-3xl md:text-5xl">Make your server a place to ship.</h2>
        <p className="mx-auto mt-5 max-w-xl text-muted-foreground">
          Install on Linux, create your account in the browser and deploy your first service. The
          installer handles the initial configuration.
        </p>
        <div className="mx-auto mt-8 max-w-xl">
          <InstallCommand />
        </div>
        <a href={productLinks.installation} className="mt-6 inline-block font-medium text-primary">
          Read the installation guide →
        </a>
      </div>
    </section>
  );
}
