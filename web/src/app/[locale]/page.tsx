import { redirect } from 'next/navigation';

export default async function LocaleRootPage({ params }: { params: Promise<{ locale: string }> }) {
  const { locale } = await params;
  // Use server-side redirect to avoid performance timing issues in the browser
  redirect(`/${locale}/admin/dashboard`);
}
