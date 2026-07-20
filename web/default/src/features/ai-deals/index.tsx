import { Link } from '@tanstack/react-router'
import {
  AlertTriangle,
  ArrowRight,
  CheckCircle2,
  CircleDollarSign,
  Globe2,
  ShieldCheck,
  Sparkles,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { PublicLayout } from '@/components/layout'
import { PageTransition } from '@/components/page-transition'

const PROVIDER_DEALS = [
  {
    provider: 'Gemini AI Pro',
    officialPrice: 'US$19.99/mo',
    localOffer: '¥69 guided subscription support / ¥129 annual reference',
    marginSignal: 'High',
    fit: 'Consumer AI membership and bundled gateway users',
    note: 'Good entry plan for price-sensitive users when delivered through official subscription channels.',
  },
  {
    provider: 'Claude Max 5x',
    officialPrice: 'US$100/mo',
    localOffer: '¥899/mo official recharge reference',
    marginSignal: 'Medium',
    fit: 'Heavy Claude users and managed workspace customers',
    note: 'Position around official recharge, subscription continuity, and assisted onboarding.',
  },
  {
    provider: 'Claude Max 20x',
    officialPrice: 'US$200/mo',
    localOffer: '¥1599-1660/mo reference',
    marginSignal: 'Low-Medium',
    fit: 'Power users who need priority Claude capacity',
    note: 'Lower gross margin; sell with SLA, official-plan guidance, and renewal health support.',
  },
  {
    provider: 'ChatGPT Plus / Pro',
    officialPrice: 'US$20+/mo',
    localOffer: 'Managed recharge quote after region check',
    marginSignal: 'Variable',
    fit: 'Users who want GPT app access plus API gateway fallback',
    note: 'Regional App Store pricing and payment eligibility require compliance-first customer qualification.',
  },
  {
    provider: 'Minimax / GLM / Mimo',
    officialPrice: 'Free tier + usage plans',
    localOffer: 'Low-cost starter bundle reference',
    marginSignal: 'Acquisition',
    fit: 'New users, testing teams, and model comparison buyers',
    note: 'Use as lead magnets and route overflow to economical models during peak demand.',
  },
  {
    provider: 'Grok / Anthropic / OpenAI API',
    officialPrice: 'Usage-based',
    localOffer: 'Gateway credits + model routing package',
    marginSignal: 'Stable',
    fit: 'Developers and small businesses using API keys',
    note: 'Bundle transparent API credits, model fallback, logs, and quota management instead of gray-market account resale.',
  },
]

const MARKET_STEPS = [
  'Acquire demand with a free comparison page, official plan list, and transparent risks.',
  'Convert users into guided subscription support, managed onboarding, or gateway-credit packages based on usage intent.',
  'Route peak traffic to cost-effective models while preserving user-facing quality and availability.',
  'Retain users with renewal reminders, subscription health checks, usage reports, and support guarantees.',
]

const SOURCE_LINKS = [
  {
    label: 'App Store regional price reference',
    href: 'https://appstoreprice.org/zh/apps',
  },
  {
    label: 'AI subscription App Store guide',
    href: 'https://appstoreprice.org/zh/blog/pay-for-ai-subscriptions-via-app-store',
  },
  { label: 'Anthropic pricing', href: 'https://www.anthropic.com/pricing' },
  { label: 'Gemini Advanced', href: 'https://gemini.google/advanced/' },
]

const SERVICE_PACKAGES = [
  {
    name: 'Starter Subscription Support',
    price: '¥69+',
    description: 'Best for users who need one consumer AI subscription with guided renewal support.',
    highlights: ['Official-channel guidance', 'Renewal reminder', 'Basic subscription health check'],
  },
  {
    name: 'Managed Premium Onboarding',
    price: '¥129-899',
    description: 'Best for customers who prefer assisted setup with documented after-sales support.',
    highlights: ['Subscription continuity support', 'Onboarding checklist', 'Support policy window'],
  },
  {
    name: 'Gateway Business Bundle',
    price: 'Custom',
    description: 'Best for teams that want API access, fallback routing, and predictable billing.',
    highlights: ['Smart model fallback', 'Quota and usage logs', 'Priority support and SLA options'],
  },
]

function SectionHeading(props: { eyebrow: string; title: string; description: string }) {
  const { t } = useTranslation()
  return (
    <div className='mx-auto max-w-3xl text-center'>
      <Badge variant='outline' className='mb-3 rounded-full px-3 py-1'>
        {t(props.eyebrow)}
      </Badge>
      <h2 className='text-2xl font-semibold tracking-tight sm:text-3xl'>
        {t(props.title)}
      </h2>
      <p className='text-muted-foreground mt-3 text-sm leading-relaxed sm:text-base'>
        {t(props.description)}
      </p>
    </div>
  )
}

export function AiDeals() {
  const { t } = useTranslation()

  return (
    <PublicLayout showMainContainer={false}>
      <div className='relative overflow-hidden'>
        <div
          aria-hidden
          className='pointer-events-none absolute inset-x-0 top-0 h-[680px] opacity-25 dark:opacity-15'
          style={{
            background: [
              'radial-gradient(ellipse 55% 42% at 20% 10%, oklch(0.72 0.18 250 / 72%) 0%, transparent 72%)',
              'radial-gradient(ellipse 48% 38% at 80% 18%, oklch(0.76 0.16 155 / 54%) 0%, transparent 70%)',
              'radial-gradient(ellipse 36% 30% at 50% 55%, oklch(0.70 0.12 40 / 42%) 0%, transparent 72%)',
            ].join(', '),
            maskImage: 'linear-gradient(to bottom, black 45%, transparent 100%)',
            WebkitMaskImage:
              'linear-gradient(to bottom, black 45%, transparent 100%)',
          }}
        />

        <PageTransition className='relative mx-auto w-full max-w-7xl px-4 pt-20 pb-14 sm:px-6 sm:pt-28 sm:pb-20 xl:px-8'>
          <header className='mx-auto max-w-4xl text-center'>
            <Badge className='mb-5 rounded-full px-3 py-1'>
              <Sparkles className='mr-1 size-3' />
              {t('AI Subscription Deal Desk')}
            </Badge>
            <h1 className='text-[clamp(2.4rem,6vw,4.8rem)] leading-[1.04] font-semibold tracking-tight'>
              {t('Compare AI subscription margins and sell compliant service bundles')}
            </h1>
            <p className='text-muted-foreground mx-auto mt-5 max-w-3xl text-base leading-relaxed sm:text-lg'>
              {t(
                'A public-facing page for AI membership price intelligence, official recharge positioning, gateway bundles, and risk-aware customer education.'
              )}
            </p>
            <div className='mt-8 flex flex-wrap items-center justify-center gap-3'>
              <Button
                className='rounded-full px-5'
                render={<Link to='/pricing' />}
              >
                {t('View Model Pricing')}
                <ArrowRight className='ml-1.5 size-4' />
              </Button>
              <Button
                variant='outline'
                className='rounded-full px-5'
                render={<a href='#business-loop' />}
              >
                {t('See Business Loop')}
              </Button>
            </div>
          </header>

          <section className='mt-12 grid gap-4 md:grid-cols-3'>
            <div className='bg-card/80 rounded-3xl border p-6 shadow-sm backdrop-blur'>
              <CircleDollarSign className='text-primary mb-4 size-8' />
              <p className='text-3xl font-semibold'>¥69+</p>
              <p className='text-muted-foreground mt-2 text-sm'>
                {t('Entry guided subscription support for price-sensitive AI membership buyers')}
              </p>
            </div>
            <div className='bg-card/80 rounded-3xl border p-6 shadow-sm backdrop-blur'>
              <Globe2 className='text-primary mb-4 size-8' />
              <p className='text-3xl font-semibold'>6+</p>
              <p className='text-muted-foreground mt-2 text-sm'>
                {t('Provider categories tracked across subscriptions and API access')}
              </p>
            </div>
            <div className='bg-card/80 rounded-3xl border p-6 shadow-sm backdrop-blur'>
              <ShieldCheck className='text-primary mb-4 size-8' />
              <p className='text-3xl font-semibold'>{t('Compliance first')}</p>
              <p className='text-muted-foreground mt-2 text-sm'>
                {t('Focus on official subscriptions, transparent support, and no abuse of trials')}
              </p>
            </div>
          </section>

          <section className='mt-20'>
            <SectionHeading
              eyebrow='Market Matrix'
              title='AI subscription and gateway offer map'
              description='Use this table as a customer-facing starting point. Pricing changes frequently, so quote final orders after manual verification.'
            />
            <div className='mt-8 overflow-hidden rounded-3xl border bg-card/70 shadow-sm backdrop-blur'>
              <div className='overflow-x-auto'>
                <table className='w-full min-w-[980px] text-left text-sm'>
                  <thead className='bg-muted/50 text-muted-foreground'>
                    <tr>
                      <th className='px-5 py-4 font-medium'>{t('Provider')}</th>
                      <th className='px-5 py-4 font-medium'>{t('Official Price')}</th>
                      <th className='px-5 py-4 font-medium'>{t('Local Offer')}</th>
                      <th className='px-5 py-4 font-medium'>{t('Margin Signal')}</th>
                      <th className='px-5 py-4 font-medium'>{t('Best Fit')}</th>
                      <th className='px-5 py-4 font-medium'>{t('Positioning Note')}</th>
                    </tr>
                  </thead>
                  <tbody className='divide-y'>
                    {PROVIDER_DEALS.map((deal) => (
                      <tr key={deal.provider} className='hover:bg-muted/30'>
                        <td className='px-5 py-4 font-semibold'>{deal.provider}</td>
                        <td className='px-5 py-4'>{deal.officialPrice}</td>
                        <td className='px-5 py-4'>{deal.localOffer}</td>
                        <td className='px-5 py-4'>
                          <Badge variant='secondary'>{t(deal.marginSignal)}</Badge>
                        </td>
                        <td className='px-5 py-4 text-muted-foreground'>{t(deal.fit)}</td>
                        <td className='px-5 py-4 text-muted-foreground'>{t(deal.note)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </section>

          <section id='business-loop' className='mt-20'>
            <SectionHeading
              eyebrow='Business Loop'
              title='A complete loop from comparison traffic to retained revenue'
              description='The safest commercial loop combines price education, official subscription support, managed onboarding, and gateway routing instead of relying on trial abuse.'
            />
            <div className='mt-8 grid gap-4 lg:grid-cols-4'>
              {MARKET_STEPS.map((step, index) => (
                <div key={step} className='bg-card/80 rounded-3xl border p-5 shadow-sm'>
                  <div className='bg-primary/10 text-primary mb-4 flex size-9 items-center justify-center rounded-full text-sm font-semibold'>
                    {index + 1}
                  </div>
                  <p className='text-sm leading-relaxed'>{t(step)}</p>
                </div>
              ))}
            </div>
          </section>

          <section className='mt-20 grid gap-6 lg:grid-cols-[0.9fr_1.1fr] lg:items-start'>
            <div className='rounded-3xl border bg-destructive/5 p-6'>
              <div className='flex items-center gap-3'>
                <AlertTriangle className='text-destructive size-6' />
                <h2 className='text-xl font-semibold'>{t('Risk boundaries')}</h2>
              </div>
              <div className='text-muted-foreground mt-4 space-y-3 text-sm leading-relaxed'>
                <p>
                  {t(
                    'Do not advertise unlimited free trials, card bypass methods, identity evasion, or guaranteed account survival. These claims create chargeback, platform, and legal risk.'
                  )}
                </p>
                <p>
                  {t(
                    'Keep customer messaging focused on official purchases, transparent fees, renewal support, onboarding education, and compliant gateway credits.'
                  )}
                </p>
              </div>
            </div>
            <div className='grid gap-4 sm:grid-cols-3'>
              {SERVICE_PACKAGES.map((item) => (
                <div key={item.name} className='bg-card rounded-3xl border p-5 shadow-sm'>
                  <Badge variant='outline' className='mb-4 rounded-full'>
                    {item.price}
                  </Badge>
                  <h3 className='font-semibold'>{t(item.name)}</h3>
                  <p className='text-muted-foreground mt-2 text-sm leading-relaxed'>
                    {t(item.description)}
                  </p>
                  <ul className='mt-4 space-y-2 text-sm'>
                    {item.highlights.map((highlight) => (
                      <li key={highlight} className='flex gap-2'>
                        <CheckCircle2 className='text-primary mt-0.5 size-4 shrink-0' />
                        <span>{t(highlight)}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              ))}
            </div>
          </section>

          <section className='mt-20 rounded-3xl border bg-card/80 p-6 shadow-sm sm:p-8'>
            <div className='flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between'>
              <div>
                <h2 className='text-xl font-semibold'>{t('Reference sources')}</h2>
                <p className='text-muted-foreground mt-2 max-w-2xl text-sm leading-relaxed'>
                  {t(
                    'Sources are provided for manual verification. Regional prices, app store policies, exchange rates, taxes, and provider rules can change without notice.'
                  )}
                </p>
              </div>
              <div className='flex flex-wrap gap-2'>
                {SOURCE_LINKS.map((source) => (
                  <Button
                    key={source.href}
                    variant='outline'
                    size='sm'
                    render={
                      <a
                        href={source.href}
                        target='_blank'
                        rel='noopener noreferrer'
                      />
                    }
                  >
                    {t(source.label)}
                  </Button>
                ))}
              </div>
            </div>
          </section>
        </PageTransition>
      </div>
    </PublicLayout>
  )
}
