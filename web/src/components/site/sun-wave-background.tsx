"use client";

import React from "react";

export default function SunWaveBackground() {
  return (
    <div className="fixed inset-0 -z-10 w-full h-full overflow-hidden bg-[#f5f5f3] dark:bg-slate-950 pointer-events-none">
      {/* The Sun: Centered, floating above the wave */}
      <div
        className="absolute top-[15%] left-1/2 -translate-x-1/2 w-[35vw] h-[35vw] md:w-[20vw] md:h-[20vw] max-w-[280px] max-h-[280px] rounded-full pointer-events-none animate-sun-breathe"
        style={{
          backgroundColor: '#e86a33',
          boxShadow: '0 0 100px 30px rgba(232,106,51,0.25), 0 0 6px rgba(232,106,51,0.5)',
        }}
      />
      {/* Dark mode sun overlay */}
      <div
        className="absolute top-[15%] left-1/2 -translate-x-1/2 w-[35vw] h-[35vw] md:w-[20vw] md:h-[20vw] max-w-[280px] max-h-[280px] rounded-full pointer-events-none opacity-0 dark:opacity-100 transition-opacity duration-500 animate-sun-breathe"
        style={{
          backgroundColor: '#6B4423',
          boxShadow: '0 0 80px 25px rgba(107,68,35,0.35)',
        }}
      />

      {/* The Wave: Smooth organic curve at the bottom */}
      <div className="absolute bottom-0 left-0 w-full leading-[0] z-10">
        <svg
          className="relative block w-[calc(100%+40px)] -ml-5 animate-wave-drift"
          style={{ height: '12vh', minHeight: '80px', maxHeight: '150px' }}
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 1440 320"
          preserveAspectRatio="none"
        >
          <path
            className="fill-slate-900 dark:fill-slate-800"
            d="M0,224L48,213.3C96,203,192,181,288,181.3C384,181,480,203,576,213.3C672,224,768,224,864,208C960,192,1056,160,1152,165.3C1248,171,1344,213,1392,234.7L1440,256L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"
          />
        </svg>
      </div>
    </div>
  );
}
