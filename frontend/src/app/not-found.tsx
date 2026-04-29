import { NextPage } from "next";

export default function NotFound() {
  return (
    <div className="container py-16 text-center">
      <h1 className="text-4xl font-bold mb-4">404</h1>
      <p className="text-muted-foreground mb-6">页面不存在</p>
      <a
        href="/"
        className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
      >
        返回首页
      </a>
    </div>
  );
}