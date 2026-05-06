"use client";

import { useEffect } from "react";

interface ErrorBoundaryProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function ErrorBoundary({ error, reset }: ErrorBoundaryProps) {
  useEffect(() => {
    console.error("Error:", error);
  }, [error]);

  return (
    <div className="container py-16 text-center">
      <h2 className="text-2xl font-bold mb-4">出错了</h2>
      <p className="text-muted-foreground mb-6">
        {error.message || "发生了未知错误"}
      </p>
      <button
        onClick={reset}
        className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
      >
        重试
      </button>
    </div>
  );
}