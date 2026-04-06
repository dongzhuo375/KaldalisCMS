"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { Button } from "@/components/ui/button";
import { Link } from '@/i18n/routing';
import { ArrowRight } from "lucide-react";
import { motion } from "framer-motion";

export default function HomePage() {
  const { isLoggedIn } = useAuthStore();

  return (
    <div className="relative min-h-[calc(100vh-4rem)] flex flex-col justify-center px-6 md:px-12 lg:px-20">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.8, ease: "easeOut" }}
        className="max-w-3xl space-y-8 text-center mx-auto"
      >
        <h1 className="text-3xl md:text-5xl lg:text-6xl font-serif font-medium tracking-tight text-foreground leading-[1.2]">
          I power experiences that{" "}
          <span className="text-accent italic underline decoration-accent/30 decoration-wavy underline-offset-4">
            spark joy.
          </span>
        </h1>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.4, duration: 0.6 }}
        >
          {isLoggedIn ? (
            <Link href="/admin/dashboard">
              <Button size="lg" className="h-12 px-8 text-base rounded-full bg-primary text-primary-foreground hover:bg-primary/90 font-bold shadow-lg shadow-primary/10 hover:shadow-xl transition-all">
                Go to Dashboard <ArrowRight className="ml-2 w-4 h-4" />
              </Button>
            </Link>
          ) : (
            <Link href="/login">
              <Button size="lg" className="h-12 px-8 text-base rounded-full bg-primary text-primary-foreground hover:bg-primary/90 font-bold shadow-lg shadow-primary/10 hover:shadow-xl transition-all">
                Get Started <ArrowRight className="ml-2 w-4 h-4" />
              </Button>
            </Link>
          )}
        </motion.div>
      </motion.div>
    </div>
  );
}
