"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {Link} from '@/i18n/routing';
import { ArrowRight, Zap, Shield, Layers } from "lucide-react";
import { useTranslations } from 'next-intl';
import SunWaveBackground from "@/components/site/sun-wave-background";
import { motion } from "framer-motion";

export default function HomePage() {
  const { isLoggedIn } = useAuthStore();
  const t = useTranslations();

  return (
    <div className="relative z-0 min-h-[calc(100vh-4rem)] flex flex-col justify-center overflow-hidden">
      <SunWaveBackground />

      <div className="relative z-10 space-y-24 py-12 md:py-20 max-w-7xl mx-auto px-6">
        
        {/* Hero 区域 */}
        <section className="text-left space-y-10 max-w-4xl">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.8 }}
            className="space-y-6"
          >
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-accent/10 border border-accent/20 text-accent text-xs font-bold tracking-widest uppercase">
              <span className="w-1.5 h-1.5 rounded-full bg-accent animate-pulse" />
              Engine v2.4.0
            </div>

            <h1 className="text-6xl md:text-8xl font-serif font-medium tracking-tight text-foreground leading-[1.1]">
              Hey! I'm Kaldalis.<br />
              I power experiences that <span className="text-accent italic">spark joy.</span>
            </h1>
            
            <p className="text-xl md:text-2xl text-muted-foreground max-w-2xl leading-relaxed font-sans">
              A minimalist headless CMS designed for clarity, speed, and creative freedom. Manage your content without the clutter.
            </p>
          </motion.div>

          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.5, duration: 0.8 }}
            className="flex flex-col sm:flex-row gap-4"
          >
              {isLoggedIn ? (
                <Link href="/admin/dashboard">
                  <Button size="lg" className="h-14 px-8 text-base rounded-full bg-primary text-primary-foreground hover:bg-primary/90 transition-all font-bold">
                    Go to Dashboard <ArrowRight className="ml-2 w-4 h-4" />
                  </Button>
                </Link>
              ) : (
                <>
                  <Link href="/login">
                    <Button size="lg" className="h-14 px-8 text-base rounded-full bg-primary text-primary-foreground hover:bg-primary/90 transition-all font-bold">
                      Get Started
                    </Button>
                  </Link>
                  <Button size="lg" variant="outline" className="h-14 px-8 text-base rounded-full border-foreground/10 hover:bg-foreground/5 font-bold">
                    View Demo
                  </Button>
                </>
              )}
          </motion.div>
        </section>

        {/* Features Grid */}
        <section className="grid grid-cols-1 md:grid-cols-3 gap-8 pb-20">
          <FeatureCard 
            icon={<Zap className="w-6 h-6 text-accent" />}
            title="Instant Speed"
            desc="Built on a high-performance Go backend for sub-millisecond response times."
          />
          <FeatureCard 
            icon={<Shield className="w-6 h-6 text-accent" />}
            title="Bulletproof RBAC"
            desc="Fine-grained access control that grows with your team and your needs."
          />
          <FeatureCard 
            icon={<Layers className="w-6 h-6 text-accent" />}
            title="Headless Freedom"
            desc="Deliver your content to any device, any platform, any time via JSON."
          />
        </section>
      </div>
    </div>
  );
}

function FeatureCard({ icon, title, desc }: { icon: React.ReactNode, title: string, desc: string }) {
  return (
    <Card className="bg-background/40 backdrop-blur-md border border-foreground/5 shadow-none hover:border-accent/20 transition-all group">
      <CardHeader className="pb-4">
        <div className="mb-2 p-2 w-fit rounded-xl bg-accent/5 group-hover:bg-accent/10 transition-colors">
          {icon}
        </div>
        <CardTitle className="text-xl font-serif font-medium text-foreground">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-muted-foreground leading-relaxed">{desc}</p>
      </CardContent>
    </Card>
  )
}
