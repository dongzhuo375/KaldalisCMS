import { redirect } from 'next/navigation';

export default function LocaleRootPage({ params }: { params: { locale: string } }) {
  // Since this is a Headless CMS Admin UI, redirect the root path to the admin dashboard.
  redirect(`/${params.locale}/admin/dashboard`);
}
