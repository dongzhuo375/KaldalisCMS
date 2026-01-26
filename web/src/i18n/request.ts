import {getRequestConfig} from 'next-intl/server';
import {routing} from './routing';
 
export default getRequestConfig(async ({requestLocale}) => {
  // This typically corresponds to the `[locale]` segment
  let locale = await requestLocale;
 
  // Ensure that a valid locale is used
  if (!locale || !routing.locales.includes(locale as any)) {
    locale = routing.defaultLocale;
  }
 
  let messages;
  try {
    switch (locale) {
      case 'en':
        messages = (await import('../messages/en.json')).default;
        break;
      case 'zh-CN':
        messages = (await import('../messages/zh-CN.json')).default;
        break;
      default:
        messages = (await import(`../messages/${locale}.json`)).default;
    }
  } catch (error) {
    // Fallback to default locale
    messages = (await import('../messages/zh-CN.json')).default;
  }
 
  return {
    locale,
    messages
  };
});