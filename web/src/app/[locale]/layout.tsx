import {NextIntlClientProvider} from 'next-intl';
import {getMessages, getTranslations} from 'next-intl/server';
import {Inter} from "next/font/google";
import "../globals.css";
import { ThemeProvider } from "@/components/theme-provider";
import { SystemStatusGuard } from "@/components/system-status-guard";
import { QueryProvider } from "@/components/providers/query-provider";
import { Toaster } from "@/components/ui/sonner";

import { redirect } from 'next/navigation';
import { routing } from '@/i18n/routing';

const inter = Inter({subsets: ["latin"]});

export async function generateMetadata() {
  const t = await getTranslations('common');

  return {
    title: "Kaldalis CMS",
    description: "A modern content management system built with Go and Next.js"
  };
}

export default async function RootLayout({
  children,
  params
}: {
  children: React.ReactNode;
  params: Promise<{locale: string}>;
}) {
  const resolvedParams = await params;
  const locale = (resolvedParams.locale && resolvedParams.locale !== 'undefined') 
    ? resolvedParams.locale 
    : routing.defaultLocale;

  // Providing all messages to the client
  // side is the easiest way to get started
  const messages = await getMessages({locale});

  return (
    <html lang={locale} suppressHydrationWarning>
      <body className={inter.className}>
        <NextIntlClientProvider locale={locale} messages={messages}>
          <QueryProvider>
            <ThemeProvider
              attribute="class"
              defaultTheme="system"
              enableSystem
              disableTransitionOnChange
            >
              <SystemStatusGuard>
                {children}
              </SystemStatusGuard>
              <Toaster position="top-right" richColors />
            </ThemeProvider>
          </QueryProvider>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}