import {NextIntlClientProvider} from 'next-intl';
import {getMessages, getTranslations} from 'next-intl/server';
import {Inter} from "next/font/google";
import "../globals.css";
import LanguageSwitcher from "@/components/LanguageSwitcher";

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
  const {locale} = await params;
  // Providing all messages to the client
  // side is the easiest way to get started
  const messages = await getMessages();

  return (
    <html lang={locale}>
      <body className={inter.className}>
        <NextIntlClientProvider messages={messages}>
          <div className="absolute top-4 right-4 z-50">
            <LanguageSwitcher />
          </div>
          {children}
        </NextIntlClientProvider>
      </body>
    </html>
  );
}