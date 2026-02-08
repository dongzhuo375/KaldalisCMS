"use client";

import { useEffect, useState } from "react";

export default function FluidBackground() {
  const [mounted, setMounted] = useState(false);
  
  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) return null;

  return (
    <div className="absolute inset-0 -z-10 w-full h-full overflow-hidden pointer-events-none select-none">
      <div className="absolute top-[-10%] left-[-10%] w-[50%] h-[50%] bg-purple-500/30 rounded-full blur-[120px] mix-blend-multiply animate-blob filter dark:bg-purple-900/20 dark:mix-blend-normal"></div>
      <div className="absolute top-[-10%] right-[-10%] w-[50%] h-[50%] bg-blue-500/30 rounded-full blur-[120px] mix-blend-multiply animate-blob animation-delay-2000 filter dark:bg-blue-900/20 dark:mix-blend-normal"></div>
      <div className="absolute bottom-[-20%] left-[20%] w-[60%] h-[60%] bg-emerald-500/30 rounded-full blur-[120px] mix-blend-multiply animate-blob animation-delay-4000 filter dark:bg-emerald-900/20 dark:mix-blend-normal"></div>
    </div>
  );
}
