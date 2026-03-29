import createNextIntlPlugin from 'next-intl/plugin';
import type { NextConfig } from "next";

const withNextIntl = createNextIntlPlugin('./src/i18n/request.ts');

const nextConfig: NextConfig = {
  /* config options here */
  async rewrites() {
    return [
      {
        source: '/api/v1/:path*',
        destination: 'http://localhost:8080/api/v1/:path*',
      },
      {
        source: '/healthz',
        destination: 'http://localhost:8080/healthz',
      },
      {
        source: '/readyz',
        destination: 'http://localhost:8080/readyz',
      },
      {
        source: '/media/:path*',
        destination: 'http://localhost:8080/media/:path*',
      },
    ]
  },
};

export default withNextIntl(nextConfig);
