"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { Button } from "@/components/ui/button";
import { Link } from '@/i18n/routing';
import { ArrowRight, Zap, Shield, Layers } from "lucide-react";
import { useTranslations } from 'next-intl';
import { motion } from "framer-motion";

export default function HomePage() {
  const { isLoggedIn } = useAuthStore();
  const t = useTranslations();

  return (
    <div className="relative">
      {/* Hero Section - Full viewport height */}
      <section className="min-h-[calc(100vh-4rem)] flex flex-col justify-center px-6 md:px-12 lg:px-20 max-w-5xl">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, ease: "easeOut" }}
          className="space-y-8"
        >
          <h1 className="text-5xl md:text-7xl lg:text-8xl font-serif font-medium tracking-tight text-foreground leading-[1.1]">
            Hey! I'm Kaldalis.<br />
            I power experiences that{" "}
            <span className="text-accent italic underline decoration-accent/30 decoration-wavy underline-offset-8">
              spark joy.
            </span>
          </h1>

          <p className="text-xl md:text-2xl text-muted-foreground max-w-2xl leading-relaxed">
            A minimalist headless CMS for modern creators.
          </p>

          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.4, duration: 0.6 }}
          >
            {isLoggedIn ? (
              <Link href="/admin/dashboard">
                <Button size="lg" className="h-14 px-10 text-base rounded-full bg-primary text-primary-foreground hover:bg-primary/90 font-bold shadow-lg shadow-primary/10 hover:shadow-xl hover:shadow-primary/20 transition-all">
                  Go to Dashboard <ArrowRight className="ml-2 w-5 h-5" />
                </Button>
              </Link>
            ) : (
              <Link href="/login">
                <Button size="lg" className="h-14 px-10 text-base rounded-full bg-primary text-primary-foreground hover:bg-primary/90 font-bold shadow-lg shadow-primary/10 hover:shadow-xl hover:shadow-primary/20 transition-all">
                  Get Started <ArrowRight className="ml-2 w-5 h-5" />
                </Button>
              </Link>
            )}
          </motion.div>
        </motion.div>
      </section>

      {/* Features Section - On the wave (dark background) */}
      <section className="relative z-20 bg-slate-900 dark:bg-slate-800 py-20 md:py-28 -mt-[12vh]">
        <div className="max-w-6xl mx-auto px-6 md:px-12">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: "-100px" }}
            transition={{ duration: 0.6 }}
            className="text-center mb-16"
          >
            <h2 className="text-3xl md:text-4xl font-serif font-medium text-white mb-4">
              Why Kaldalis?
            </h2>
            <p className="text-slate-400 text-lg max-w-xl mx-auto">
              Three things that make us different.
            </p>
          </motion.div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <FeatureCard
              icon={<Zap className="w-7 h-7" />}
              title="Instant Speed"
              desc="Built on a high-performance Go backend for sub-millisecond response times."
              delay={0}
            />
            <FeatureCard
              icon={<Shield className="w-7 h-7" />}
              title="Bulletproof RBAC"
              desc="Fine-grained access control that grows with your team and your needs."
              delay={0.1}
            />
            <FeatureCard
              icon={<Layers className="w-7 h-7" />}
              title="Headless Freedom"
              desc="Deliver your content to any device, any platform, any time via JSON."
              delay={0.2}
            />
          </div>
        </div>
      </section>
    </div>
  );
}

function FeatureCard({
  icon,
  title,
  desc,
  delay = 0
}: {
  icon: React.ReactNode;
  title: string;
  desc: string;
  delay?: number;
}) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, margin: "-50px" }}
      transition={{ duration: 0.5, delay }}
      className="p-8 rounded-2xl bg-slate-800/50 border border-slate-700/50 hover:border-accent/30 transition-colors group"
    >
      <div className="mb-5 p-3 w-fit rounded-xl bg-accent/10 text-accent group-hover:bg-accent group-hover:text-white transition-colors">
        {icon}
      </div>
      <h3 className="text-xl font-serif font-medium text-white mb-3">{title}</h3>
      <p className="text-slate-400 leading-relaxed">{desc}</p>
    </motion.div>
  );
}
