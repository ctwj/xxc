export default function Loading() {
  return (
    <div className="container py-16">
      <div className="flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
      <p className="text-center text-muted-foreground mt-4">加载中...</p>
    </div>
  );
}