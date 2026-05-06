"use client";

import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { InfoCard } from './InfoCard';
import { InfoCardData } from '@/types';
import { RefreshCw } from 'lucide-react';

interface CardStackProps {
  cards: InfoCardData[];
  currentIndex: number;
  onSwipe: (id: string, direction: 'left' | 'right') => void;
  onFavorite: (card: InfoCardData) => void;
  isFavorited: (id: string) => boolean;
  onRefresh?: () => void;
  onShare?: (card: InfoCardData) => void;
}

export const CardStack: React.FC<CardStackProps> = ({
  cards,
  currentIndex,
  onSwipe,
  onFavorite,
  isFavorited,
  onRefresh,
  onShare
}) => {
  const isLastCard = currentIndex >= cards.length;

  if (isLastCard || cards.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center w-full h-full p-8 text-center bg-background">
        <motion.div
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          className="max-w-md"
        >
          <div className="w-24 h-24 bg-card rounded-full flex items-center justify-center mx-auto mb-8 shadow-inner">
            <RefreshCw size={40} className="text-primary/40" />
          </div>
          <h2 className="text-3xl font-bold mb-4">今日已看完</h2>
          <p className="text-foreground/60 mb-10 leading-relaxed">
            你已经浏览完所有的精选内容了，明天再来看看吧，或者重新开始。
          </p>
          <button
            onClick={() => {
              onRefresh?.();
            }}
            className="px-10 py-4 bg-primary text-white rounded-full font-bold shadow-lg shadow-primary/30 active:scale-95 transition-all flex items-center gap-3 mx-auto"
          >
            <RefreshCw size={20} />
            重新开始
          </button>
        </motion.div>
      </div>
    );
  }

  return (
    <div className="relative w-full h-full flex items-center justify-center overflow-hidden perspective-1000">
      {/* 只显示当前卡片 */}
      {cards[currentIndex] && (
        <motion.div
          key={`card-${currentIndex}`}
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          exit={{ opacity: 0, scale: 0.9 }}
          transition={{ duration: 0.3 }}
          className="absolute w-[90%] md:w-[450px] aspect-[9/16] md:h-[80vh]"
        >
          <InfoCard
            data={cards[currentIndex]}
            isCurrent={true}
            onShare={onShare ? () => onShare(cards[currentIndex]) : undefined}
            className="h-full w-full"
          />
        </motion.div>
      )}
    </div>
  );
};