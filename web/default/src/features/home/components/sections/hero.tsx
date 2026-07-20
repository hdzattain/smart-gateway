/*
Copyright (C) 2023-2026 Smart Gateway

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@smart-gateway.shop
*/
import { Link } from '@tanstack/react-router'
import { ArrowRight, BookOpen } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { useStatus } from '@/hooks/use-status'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { HeroTerminalDemo } from '../hero-terminal-demo'

interface HeroProps {
  className?: string
  isAuthenticated?: boolean
}

const TRUSTED_MODELS = ['OpenAI', 'Claude', 'Gemini', 'Azure', 'Bedrock']

export function Hero(props: HeroProps) {
  const { t } = useTranslation()
  const { status } = useStatus()
  const docsUrl =
    (status?.docs_link as string | undefined) || 'https://smart-gateway.shop/docs'

  const renderDocsButton = () => {
    const isExternal = docsUrl.startsWith('http')
    if (isExternal) {
      return (
        <Button
          variant='outline'
          className='group border-border/50 hover:border-border hover:bg-muted/50 inline-flex h-11 items-center gap-1.5 rounded-lg px-5 text-sm font-medium'
          render={
            <a href={docsUrl} target='_blank' rel='noopener noreferrer' />
          }
        >
          <BookOpen className='text-muted-foreground/80 group-hover:text-foreground size-4 transition-colors duration-200' />
          <span>{t('Docs')}</span>
        </Button>
      )
    }
    return (
      <Button
        variant='outline'
        className='group border-border/50 hover:border-border hover:bg-muted/50 inline-flex h-11 items-center gap-1.5 rounded-lg px-5 text-sm font-medium'
        render={<Link to={docsUrl} />}
      >
        <BookOpen className='text-muted-foreground/80 group-hover:text-foreground size-4 transition-colors duration-200' />
        <span>{t('Docs')}</span>
      </Button>
    )
  }

  return (
    <section className='relative z-10 overflow-hidden px-6 pt-24 pb-16 md:pt-32 md:pb-24 lg:pt-36 lg:pb-28'>
      <div
        aria-hidden
        className='pointer-events-none absolute inset-x-0 top-0 -z-10 h-[42rem] bg-[radial-gradient(ellipse_72%_58%_at_50%_0%,color-mix(in_oklch,var(--primary)_12%,transparent)_0%,transparent_72%),linear-gradient(to_bottom,color-mix(in_oklch,var(--background)_96%,var(--primary)_4%)_0%,transparent_32%)] opacity-100'
      />
      <div
        aria-hidden
        className='pointer-events-none absolute inset-0 -z-10 bg-[linear-gradient(to_right,var(--border)_1px,transparent_1px),linear-gradient(to_bottom,var(--border)_1px,transparent_1px)] bg-[size:5rem_5rem] opacity-[0.05] [mask-image:linear-gradient(to_bottom,black_26%,transparent_100%)]'
      />

      <div className='mx-auto grid max-w-6xl gap-12 lg:grid-cols-[minmax(0,1.04fr)_minmax(0,0.96fr)] lg:items-center lg:gap-10'>
        <div className='flex flex-col items-start text-left'>
          <div
            className='landing-animate-fade-up mb-5 inline-flex items-center gap-2 rounded-full border border-border/60 bg-background/80 px-3 py-1.5 text-[11px] font-medium text-foreground/70 opacity-0 shadow-xs backdrop-blur-sm'
            style={{ animationDelay: '0ms' }}
          >
            <span className='size-2 rounded-full bg-emerald-500' />
            <span>{t('AI Application Infrastructure Foundation')}</span>
          </div>

          <h1
            className='landing-animate-fade-up text-[clamp(2.5rem,5vw,4.3rem)] leading-[1.02] font-semibold tracking-tight opacity-0'
            style={{ animationDelay: '60ms' }}
          >
            {t('Unified API Gateway for')}
            <br />
            <span className='text-primary'>
              {t('Vast Range of AI Models')}
            </span>
          </h1>
          <p
            className='landing-animate-fade-up text-muted-foreground/80 mt-5 max-w-xl text-[15px] leading-relaxed opacity-0 md:text-base'
            style={{ animationDelay: '120ms' }}
          >
            {t(
              'Access a vast selection of models via a standard, unified API protocol. Power AI applications, manage digital assets, and connect the Future.'
            )}
          </p>

          <div
            className='landing-animate-fade-up mt-8 flex flex-wrap items-center gap-3 opacity-0'
            style={{ animationDelay: '180ms' }}
          >
            {props.isAuthenticated ? (
              <>
                <Button
                  className='group h-11 rounded-full px-5 text-sm font-medium'
                  render={<Link to='/dashboard' />}
                >
                  {t('Go to Dashboard')}
                  <ArrowRight className='ml-1.5 size-4 transition-transform duration-200 group-hover:translate-x-0.5' />
                </Button>
                {renderDocsButton()}
              </>
            ) : (
              <>
                <Button
                  className='group h-11 rounded-full px-5 text-sm font-medium'
                  render={<Link to='/sign-up' />}
                >
                  {t('Get Started')}
                  <ArrowRight className='ml-1.5 size-4 transition-transform duration-200 group-hover:translate-x-0.5' />
                </Button>
                <Button
                  variant='outline'
                  className='border-border/60 hover:border-border hover:bg-muted/40 h-11 rounded-full px-5 text-sm font-medium'
                  render={<Link to='/pricing' />}
                >
                  {t('View Pricing')}
                </Button>
                {renderDocsButton()}
              </>
            )}
          </div>

          <div
            className='landing-animate-fade-up mt-8 flex w-full max-w-xl flex-col gap-4 opacity-0'
            style={{ animationDelay: '240ms' }}
          >
            <div className='flex flex-wrap items-center gap-2.5'>
              {TRUSTED_MODELS.map((model) => (
                <Badge
                  key={model}
                  variant='secondary'
                  className='rounded-full border border-border/60 bg-background/75 px-3 py-1 text-[11px] font-medium shadow-none'
                >
                  {model}
                </Badge>
              ))}
            </div>

            <div className='grid gap-3 sm:grid-cols-3'>
              <div className='rounded-2xl border border-border/60 bg-background/75 px-4 py-4 shadow-xs backdrop-blur-sm'>
                <p className='text-muted-foreground text-[11px] font-medium tracking-[0.16em] uppercase'>
                  {t('Route')}
                </p>
                <p className='mt-2 text-sm leading-relaxed'>
                  {t('Route, auth, and balance check in one place')}
                </p>
              </div>
              <div className='rounded-2xl border border-border/60 bg-background/75 px-4 py-4 shadow-xs backdrop-blur-sm'>
                <p className='text-muted-foreground text-[11px] font-medium tracking-[0.16em] uppercase'>
                  {t('Developer Friendly')}
                </p>
                <p className='mt-2 text-sm leading-relaxed'>
                  {t('Compatible API routes for common AI application workflows')}
                </p>
              </div>
              <div className='rounded-2xl border border-border/60 bg-background/75 px-4 py-4 shadow-xs backdrop-blur-sm'>
                <p className='text-muted-foreground text-[11px] font-medium tracking-[0.16em] uppercase'>
                  {t('Transparent Billing')}
                </p>
                <p className='mt-2 text-sm leading-relaxed'>
                  {t('Pay-as-you-go with real-time usage monitoring')}
                </p>
              </div>
            </div>
          </div>
        </div>

        <div
          className='landing-animate-fade-up flex w-full justify-center opacity-0 lg:pl-4'
          style={{ animationDelay: '320ms' }}
        >
          <div className='w-full max-w-[36rem]'>
            <div className='mb-4 flex items-center gap-2 text-[11px] font-medium tracking-[0.16em] uppercase text-muted-foreground/60'>
              <span className='size-2 rounded-full bg-primary/70' />
              <span>{t('Route active')}</span>
            </div>
            <HeroTerminalDemo />
          </div>
        </div>
      </div>
    </section>
  )
}
