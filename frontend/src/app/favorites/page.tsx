"use client";

import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { ChevronLeft, User, Bookmark, Trash2 } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { InfoCard } from '@/components/home/InfoCard';
import { InfoCardData } from '@/types';
import { cn } from '@/lib/utils';

// Mock data for favorites (will be replaced with API calls)
const MOCK_FAVORITES: InfoCardData[] = [
  {
    id: '1',
    type: 'image-text',
    title: '示例收藏文章',
    content: '这是一篇收藏的文章内容示例。',
    mediaUrl: 'https://images.unsplash.com/photo-1557683316-973673baf926?q=80&w=1000&auto=format&fit=crop',
    tag: '技术',
    publishTime: new Date().toISOString(),
  },
];

export default function FavoritePage() {
  const router = useRouter();
  const [favorites, setFavorites] = useState<InfoCardData[]>([]);
  const [selectedCard, setSelectedCard] = useState<InfoCardData | null>(null);

  useEffect(() => {
    // Load favorites from localStorage for now
    const savedFavorites = localStorage.getItem('favorites');
    if (savedFavorites) {
      setFavorites(JSON.parse(savedFavorites));
    } else {
      setFavorites(MOCK_FAVORITES);
    }
  }, []);

  const toggleFavorite = (card: InfoCardData) => {
    const newFavorites = favorites.filter(f => f.id !== card.id);
    setFavorites(newFavorites);
    localStorage.setItem('favorites', JSON.stringify(newFavorites));
  };

  return (
    <div className="min-h-screen bg-background text-foreground pb-20">
      {/* 顶部导航 */}
      <header className="fixed top-0 left-0 right-0 z-50 glass px-6 py-4 flex items-center justify-between">
        <button
          onClick={() => router.push('/')}
          className="p-2 hover:bg-white/10 rounded-full transition-colors"
        >
          <ChevronLeft size={24} />
        </button>
        <h1 className="text-lg font-bold">我的收藏</h1>
        <div className="w-10 h-10 rounded-full bg-card flex items-center justify-center border border-white/10">
          <User size={20} className="text-white/40" />
        </div>
      </header>

      {/* 用户信息 */}
      <div className="pt-32 px-8 mb-12 flex items-center gap-6">
        <div className="w-20 h-20 rounded-full bg-gradient-to-br from-primary to-primary/40 flex items-center justify-center shadow-lg shadow-primary/20">
          <User size={40} className="text-white" />
        </div>
        <div>
          <h2 className="text-2xl font-bold mb-1">每日信息差用户</h2>
          <div className="flex items-center gap-2 text-foreground/40 text-sm">
            <Bookmark size={14} />
            <span>{favorites.length} 个收藏内容</span>
          </div>
        </div>
      </div>

      {/* 收藏列表 */}
      <div className="px-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <AnimatePresence mode="popLayout">
          {favorites.length > 0 ? (
            favorites.map((card) => (
              <motion.div
                key={card.id}
                layout
                initial={{ opacity: 0, scale: 0.9 }}
                animate={{ opacity: 1, scale: 1 }}
                exit={{ opacity: 0, scale: 0.8 }}
                className="group relative h-80 bg-card rounded-[20px] overflow-hidden border border-white/5 cursor-pointer shadow-lg hover:shadow-2xl hover:border-white/10 transition-all duration-300"
                onClick={() => setSelectedCard(card)}
              >
                {/* 缩略图内容 */}
                {card.mediaUrl ? (
                  <div className="h-full w-full">
                    {card.type.includes('video') ? (
                      <div className="relative h-full w-full">
                        <img src={card.coverUrl} className="h-full w-full object-cover brightness-50" />
                        <div className="absolute inset-0 flex items-center justify-center">
                          <div className="w-12 h-12 bg-white/20 backdrop-blur-md rounded-full flex items-center justify-center">
                            <div className="w-0 h-0 border-t-[8px] border-t-transparent border-l-[14px] border-l-white border-b-[8px] border-b-transparent ml-1" />
                          </div>
                        </div>
                      </div>
                    ) : (
                      <img src={card.mediaUrl} className="h-full w-full object-cover brightness-75 group-hover:scale-110 transition-transform duration-700" />
                    )}
                  </div>
                ) : (
                  <div className="h-full w-full p-8 flex items-center justify-center text-center bg-gradient-to-br from-card to-background">
                    <p className="line-clamp-4 font-medium text-lg text-white/80">{card.content}</p>
                  </div>
                )}

                {/* 文字覆盖 */}
                <div className="absolute inset-0 bg-gradient-to-t from-black/90 via-black/20 to-transparent p-6 flex flex-col justify-end">
                  <div className="flex items-center gap-2 mb-2">
                    <span className="px-2 py-0.5 bg-primary/20 text-primary text-[10px] font-bold rounded uppercase tracking-wider">{card.tag}</span>
                  </div>
                  <h3 className="text-white font-bold leading-tight line-clamp-2">{card.title}</h3>
                </div>

                {/* 悬浮操作 */}
                <div className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      toggleFavorite(card);
                    }}
                    className="p-3 bg-destructive/20 backdrop-blur-md text-destructive hover:bg-destructive hover:text-white rounded-full transition-all"
                  >
                    <Trash2 size={18} />
                  </button>
                </div>
              </motion.div>
            ))
          ) : (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              className="col-span-full py-20 text-center"
            >
              <div className="w-24 h-24 bg-card rounded-full flex items-center justify-center mx-auto mb-6">
                <Bookmark size={40} className="text-white/20" />
              </div>
              <h3 className="text-xl font-bold mb-3">还没有收藏内容</h3>
              <p className="text-foreground/40 mb-8 max-w-xs mx-auto leading-relaxed">去首页发现感兴趣的信息并点击星星图标进行收藏吧</p>
              <button
                onClick={() => router.push('/')}
                className="rounded-full px-8 py-6 font-bold bg-primary text-white hover:bg-primary/90 transition-all"
              >
                返回首页发现
              </button>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* 详情弹窗 */}
      <AnimatePresence>
        {selectedCard && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-[100] flex items-center justify-center p-4 md:p-8"
          >
            <div
              className="absolute inset-0 bg-black/90 backdrop-blur-lg"
              onClick={() => setSelectedCard(null)}
            />
            <motion.div
              layoutId={selectedCard.id}
              initial={{ scale: 0.9, y: 50 }}
              animate={{ scale: 1, y: 0 }}
              exit={{ scale: 0.9, y: 50 }}
              className="relative w-full max-w-lg aspect-[9/16] md:h-[85vh] rounded-[20px] overflow-hidden"
            >
              <InfoCard data={selectedCard} isCurrent={true} className="h-full w-full" />
              <button
                onClick={() => setSelectedCard(null)}
                className="absolute top-8 right-8 p-3 bg-black/30 backdrop-blur-md rounded-full text-white/80 hover:bg-black/50 transition-colors z-[110]"
              >
                <ChevronLeft size={24} className="rotate-[-90deg]" />
              </button>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}