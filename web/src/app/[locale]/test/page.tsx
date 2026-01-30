"use client";

import { useTranslations } from 'next-intl';

export default function TestPage() {
  const t = useTranslations();
  
  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-4">Translation Test</h1>
      <div className="space-y-2">
        <p>App Name: {t('common.app_name')}</p>
        <p>Welcome: {t('common.welcome')}</p>
        <p>Login: {t('common.login')}</p>
        <p>Home: {t('navigation.home')}</p>
        <p>Hero Subtitle: {t('home.hero_subtitle')}</p>
        <p>Login Title: {t('auth.login_title')}</p>
      </div>
    </div>
  );
}