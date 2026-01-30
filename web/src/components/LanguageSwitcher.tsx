"use client";

import {useLocale, useTranslations} from 'next-intl';
import {usePathname, Link} from '@/i18n/routing';
import {Button} from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {Globe} from 'lucide-react';

export default function LanguageSwitcher() {
  const locale = useLocale();
  const pathname = usePathname();
  const t = useTranslations('navigation');

  const languages = [
    {code: 'zh-CN', name: '简体中文', flag: '🇨🇳'},
    {code: 'en', name: 'English', flag: '🇺🇸'},
  ];

  const currentLanguage = languages.find(lang => lang.code === locale);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="gap-2">
          <Globe className="h-4 w-4" />
          {currentLanguage?.flag} {currentLanguage?.name}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {languages.map((language) => (
          <DropdownMenuItem key={language.code} asChild>
            <Link 
              href={pathname} 
              locale={language.code}
              className={locale === language.code ? 'bg-accent w-full cursor-pointer' : 'w-full cursor-pointer'}
            >
              <span className="mr-2">{language.flag}</span>
              {language.name}
            </Link>
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}