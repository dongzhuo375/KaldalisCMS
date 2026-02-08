"use client"

import * as React from "react"
import { Check } from "lucide-react"
import { cn } from "@/lib/utils"

interface CheckboxProps extends React.InputHTMLAttributes<HTMLInputElement> {
  onCheckedChange?: (checked: boolean) => void;
}

const Checkbox = React.forwardRef<HTMLInputElement, CheckboxProps>(
  ({ className, onCheckedChange, onChange, ...props }, ref) => {
    
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      onChange?.(e);
      onCheckedChange?.(e.target.checked);
    };

    return (
      <div className="relative flex items-center">
        <input
          type="checkbox"
          className={cn(
            "peer h-4 w-4 shrink-0 rounded-sm border border-slate-200 border-slate-900 shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-slate-950 disabled:cursor-not-allowed disabled:opacity-50 dark:border-slate-800 dark:border-slate-50 dark:focus-visible:ring-slate-300 appearance-none checked:bg-slate-900 dark:checked:bg-slate-50 checked:border-transparent transition-colors cursor-pointer",
            className
          )}
          onChange={handleChange}
          ref={ref}
          {...props}
        />
        <Check className="absolute left-0 top-0 h-4 w-4 hidden peer-checked:block text-slate-50 dark:text-slate-900 pointer-events-none p-0.5" />
      </div>
    )
  }
)
Checkbox.displayName = "Checkbox"

export { Checkbox }
