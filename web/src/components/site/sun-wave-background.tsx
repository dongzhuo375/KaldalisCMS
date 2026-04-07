"use client";

import React from "react";

export default function SunWaveBackground() {
  return (
    <div className="fixed inset-0 -z-10 w-full h-full overflow-hidden bg-[var(--background)] pointer-events-none">
      {/* Sun — warm orange circle */}
      <div
        className="absolute top-[12%] left-1/2 -translate-x-1/2 w-[30vw] h-[30vw] md:w-[18vw] md:h-[18vw] max-w-[260px] max-h-[260px] rounded-full animate-sun-breathe"
        style={{
          backgroundColor: "oklch(0.6 0.25 35)",
          boxShadow:
            "0 0 120px 40px oklch(0.6 0.25 35 / 0.2), 0 0 40px 10px oklch(0.6 0.25 35 / 0.3)",
        }}
      />
      {/* Dark mode sun overlay */}
      <div
        className="absolute top-[12%] left-1/2 -translate-x-1/2 w-[30vw] h-[30vw] md:w-[18vw] md:h-[18vw] max-w-[260px] max-h-[260px] rounded-full opacity-0 dark:opacity-100 transition-opacity duration-500 animate-sun-breathe"
        style={{
          backgroundColor: "oklch(0.4 0.12 35)",
          boxShadow: "0 0 80px 25px oklch(0.4 0.12 35 / 0.35)",
        }}
      />

      {/* Organic multi-layer waves */}
      <div className="absolute bottom-0 left-0 w-full h-[30vh] min-h-[160px] max-h-[300px]">
        {/* Layer 1 — back, slowest, most transparent, tallest */}
        <svg
          className="absolute bottom-0 left-0 w-[200%] h-full animate-wave-slow"
          viewBox="0 0 1440 200"
          preserveAspectRatio="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            className="fill-foreground/[0.06] dark:fill-foreground/[0.04]"
            d="M0,120 C120,140 240,90 360,110 C480,130 540,80 720,100 C900,120 960,70 1080,95 C1200,120 1320,85 1440,105 L1440,200 L0,200Z"
          />
        </svg>

        {/* Layer 2 — middle speed, medium opacity */}
        <svg
          className="absolute bottom-0 left-0 w-[200%] h-[85%] animate-wave-mid"
          viewBox="0 0 1440 200"
          preserveAspectRatio="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            className="fill-foreground/[0.12] dark:fill-foreground/[0.08]"
            d="M0,130 C180,100 300,150 480,120 C660,90 720,140 900,125 C1080,110 1200,145 1440,115 L1440,200 L0,200Z"
          />
        </svg>

        {/* Layer 3 — front, fastest, opaque, the main visible wave matching the design reference */}
        <svg
          className="absolute bottom-0 left-0 w-[200%] h-[70%] animate-wave-fast"
          viewBox="0 0 1440 200"
          preserveAspectRatio="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            className="fill-slate-900 dark:fill-slate-800"
            d="M0,100 C60,95 120,115 200,105 C280,95 340,125 440,110 C540,95 600,130 720,115 C840,100 900,135 1020,120 C1140,105 1200,130 1320,118 C1380,112 1420,125 1440,120 L1440,200 L0,200Z"
          />
        </svg>
      </div>
    </div>
  );
}
