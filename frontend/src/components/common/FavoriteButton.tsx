"use client";

import { useState } from "react";
import { favoriteService } from "@/services/favoriteService";

interface FavoriteButtonProps {
  articleId: number;
  initialFavorited?: boolean;
  onToggle?: (favorited: boolean) => void;
}

export function FavoriteButton({
  articleId,
  initialFavorited = false,
  onToggle,
}: FavoriteButtonProps) {
  const [favorited, setFavorited] = useState(initialFavorited);
  const [loading, setLoading] = useState(false);

  const handleToggle = async () => {
    setLoading(true);
    try {
      if (favorited) {
        // Need to find the favorite ID first
        const { data } = await favoriteService.list();
        const favorite = data.find((item) => item.articleId === articleId);
        if (favorite) {
          await favoriteService.remove(favorite.id);
        }
        setFavorited(false);
        onToggle?.(false);
      } else {
        await favoriteService.add(articleId);
        setFavorited(true);
        onToggle?.(true);
      }
    } catch (error) {
      console.error("Failed to toggle favorite:", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <button
      onClick={handleToggle}
      disabled={loading}
      className={`inline-flex items-center gap-1 rounded-md px-3 py-1.5 text-sm font-medium transition-colors ${
        favorited
          ? "bg-red-100 text-red-600 hover:bg-red-200 dark:bg-red-900/20 dark:text-red-400"
          : "bg-secondary text-secondary-foreground hover:bg-secondary/80"
      } disabled:opacity-50`}
      title={favorited ? "取消收藏" : "添加收藏"}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill={favorited ? "currentColor" : "none"}
        stroke="currentColor"
        strokeWidth={favorited ? 0 : 2}
        className="h-4 w-4"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z"
        />
      </svg>
      {favorited ? "已收藏" : "收藏"}
    </button>
  );
}