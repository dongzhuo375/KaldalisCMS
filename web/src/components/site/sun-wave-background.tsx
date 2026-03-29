"use client";

import React from "react";
import { motion } from "framer-motion";

export default function SunWaveBackground() {
  return (
    <div className="fixed inset-0 -z-10 w-full h-full overflow-hidden bg-background pointer-events-none">
      {/* Subtle Grain Overlay */}
      <div className="absolute inset-0 bg-grain pointer-events-none" />

      {/* The Sun: Smaller, subtle, floating */}
      <motion.div 
        animate={{ 
          y: [0, -20, 0],
          opacity: [0.2, 0.3, 0.2]
        }}
        transition={{ 
          duration: 8, 
          repeat: Infinity, 
          ease: "easeInOut" 
        }}
        className="absolute top-[10%] left-[15%] w-[30vw] h-[30vw] max-w-[400px] max-h-[400px] rounded-full bg-accent blur-[120px] pointer-events-none z-0"
      />

      {/* The Wave: More discrete at the very bottom */}
      <div className="absolute bottom-0 left-0 w-full leading-[0] z-0 opacity-10">
        <svg 
          className="relative block w-full h-[60px] md:h-[80px]" 
          xmlns="http://www.w3.org/2000/svg" 
          viewBox="0 24 150 28" 
          preserveAspectRatio="none" 
          shapeRendering="auto"
        >
          <defs>
            <path id="gentle-wave" d="M-160 44c30 0 58-18 88-18s 58 18 88 18 58-18 88-18 58 18 88 18 v44h-352z" />
          </defs>
          <g className="parallax">
            <use xlinkHref="#gentle-wave" x="48" y="0" className="fill-foreground/[0.03] animate-[wave_25s_cubic-bezier(.55,.5,.45,.5)_infinite]" />
            <use xlinkHref="#gentle-wave" x="48" y="3" className="fill-foreground/[0.05] animate-[wave_10s_cubic-bezier(.55,.5,.45,.5)_infinite]" />
            <use xlinkHref="#gentle-wave" x="48" y="5" className="fill-foreground/[0.08] animate-[wave_15s_cubic-bezier(.55,.5,.45,.5)_infinite]" />
            <use xlinkHref="#gentle-wave" x="48" y="7" className="fill-foreground/[0.12] animate-[wave_20s_cubic-bezier(.55,.5,.45,.5)_infinite]" />
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
