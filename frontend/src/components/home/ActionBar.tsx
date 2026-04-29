"use client";

import React from 'react';
import { ThumbsDown, ThumbsUp, Star } from 'lucide-react';
import { cn } from '@/lib/utils';
import { motion } from 'framer-motion';

interface ActionBarProps {
  onDislike: () => void;
  onLike: () => void;
  onFavorite: () => void;
  isFavorited: boolean;
  className?: string;
}

export const ActionBar: React.FC<ActionBarProps> = ({
  onDislike,
  onLike,
  onFavorite,
  isFavorited,
  className
}) => {
  return (
    <div className={cn(
      "fixed bottom-12 left-1/2 -translate-x-1/2 flex items-center gap-12 z-50 glass px-12 py-5 rounded-full",
      className
    )}>
      {/* 踩按钮 */}
      <button
        onClick={onDislike}
        className="group relative flex items-center justify-center w-14 h-14 rounded-full border-2 border-destructive text-destructive hover:bg-destructive hover:text-white transition-all duration-300"
      >
        <ThumbsDown size={28} />
      </button>

      {/* 收藏按钮 */}
      <motion.button
        whileTap={{ scale: 1.4 }}
        onClick={onFavorite}
        className={cn(
          "relative flex items-center justify-center w-12 h-12 rounded-full border-2 transition-all duration-300",
          isFavorited
            ? "border-accent bg-accent text-accent-foreground"
            : "border-accent text-accent hover:bg-accent/10"
        )}
      >
        <Star size={24} fill={isFavorited ? "currentColor" : "none"} />
      </motion.button>

      {/* 赞按钮 */}
      <button
        onClick={onLike}
        className="group relative flex items-center justify-center w-14 h-14 rounded-full border-2 border-success text-success hover:bg-success hover:text-white transition-all duration-300"
      >
        <ThumbsUp size={28} />
      </button>
    </div>
  );
};