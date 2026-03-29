"use client";

import React from "react";
import { motion } from "framer-motion";

export default function SunWaveBackground() {
  return (
    <div className="fixed inset-0 -z-10 w-full h-full overflow-hidden bg-background pointer-events-none">
      {/* Subtle Grain Overlay */}
      <div className="absolute inset-0 bg-grain pointer-events-none" />

      {/* The Sun: Large blurred orange circle */}
      <motion.div 
        initial={{ scale: 0.8, opacity: 0 }}
        animate={{ scale: 1, opacity: 0.6 }}
        transition={{ duration: 2, ease: "easeOut" }}
        className="absolute -top-[10%] -left-[5%] w-[50vw] h-[50vw] max-w-[600px] max-h-[600px] rounded-full bg-accent blur-[120px] mix-blend-multiply dark:mix-blend-normal opacity-60"
      />

      {/* The Wave: Animated SVG at the bottom */}
      <div className="absolute bottom-0 left-0 w-full leading-[0]">
        <svg 
          className="relative block w-full h-[150px] md:h-[200px]" 
          xmlns="http://www.w3.org/2000/svg" 
          viewBox="0 24 150 28" 
          preserveAspectRatio="none" 
          shapeRendering="auto"
        >
          <defs>
            <path id="gentle-wave" d="M-160 44c30 0 58-18 88-18s 58 18 88 18 58-18 88-18 58 18 88 18 v44h-352z" />
          </defs>
          <g className="parallax">
            <use xlinkHref="#gentle-wave" x="48" y="0" className="fill-foreground/5 animate-[wave_25s_cubic-bezier(.55,.5,.45,.5)_infinite]" />
            <use xlinkHref="#gentle-wave" x="48" y="3" className="fill-foreground/10 animate-[wave_10s_cubic-bezier(.55,.5,.45,.5)_infinite]" />
            <use xlinkHref="#gentle-wave" x="48" y="5" className="fill-foreground/15 animate-[wave_15s_cubic-bezier(.55,.5,.45,.5)_infinite]" />
            <use xlinkHref="#gentle-wave" x="48" y="7" className="fill-foreground/20 animate-[wave_20s_cubic-bezier(.55,.5,.45,.5)_infinite]" />
          </g>
        </svg>
      </div>

      <style jsx>{`
        @keyframes wave {
          0% { transform: translate3d(-90px, 0, 0); }
          100% { transform: translate3d(85px, 0, 0); }
        }
      `}</style>
    </div>
  );
}
