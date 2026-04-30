"use client";

import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { ChevronLeft, User, Bookmark, Trash2, Loader2 } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { InfoCard } from '@/components/home/InfoCard';
import { InfoCardData, Article } from '@/types';
import { api } from '@/lib/api';
import { toast } from 'sonner';

// Convert Article to InfoCardData
function articleToCardData(article: Article, tag?: string): InfoCardData {
  // Parse mediaUrls if it's a JSON string
  let mediaUrls: string[] = [];
  if (article.mediaUrls) {
    try {
      mediaUrls = JSON.parse(article.mediaUrls);
    } catch {
      // If not valid JSON, treat as single URL
      mediaUrls = article.mediaUrls ? [article.mediaUrls] : [];
    }
  }

  // Determine content type
  let type: InfoCardData['type'] = (article.type as InfoCardData['type']) || 'text';
  if (!article.type) {
    // Auto-detect type based on available fields
    if (article.videoUrl) {
      type = article.thumbnail ? 'video-text' : 'video';
    } else if (mediaUrls.length > 1) {
      type = article.content ? 'images-text' : 'images';
    } else if (article.thumbnail) {
      type = article.content ? 'image-text' : 'image';
    } else if (article.content && article.content.length > 500) {
      type = 'long-text';
    }
  }

  return {
    id: String(article.id),
    type,
    title: article.title,
    content: article.content || article.description,
    mediaUrl: article.thumbnail,
    mediaUrls,
    videoUrl: article.videoUrl,
    coverUrl: article.coverUrl,
    tag: tag || '未分类',
    publishTime: new Date(article.createTime * 1000).toISOString(),
    isNew: Date.now() - article.createTime * 1000 < 24 * 60 * 60 * 1000,
  };
}

interface FavoriteItem {
  id: number;
  articleId: number;
  article: Article;
}

export default function FavoritePage() {
  const router = useRouter();
  const [favorites, setFavorites] = useState<InfoCardData[]>([]);
  const [favoriteIds, setFavoriteIds] = useState<Map<string, number>>(new Map()); // articleId -> favoriteId
  const [selectedCard, setSelectedCard] = useState<InfoCardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState<{ id: number; username: string } | null>(null);

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        // Check if user is logged in
        const userData = await api.getCurrentUser();
        setUser(userData);

        // Load favorites from API
        const response = await api.getFavorites();
        const favItems = response.data as FavoriteItem[];

        const cardData: InfoCardData[] = [];
        const idMap = new Map<string, number>();

        for (const fav of favItems) {
          if (fav.article) {
            cardData.push(articleToCardData(fav.article));
            idMap.set(String(fav.articleId), fav.id);
          }
        }

        setFavorites(cardData);
        setFavoriteIds(idMap);
      } catch (error) {
        // Not logged in or error
        console.error('Failed to load favorites:', error);
        setUser(null);
        router.push('/login');
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [router]);

  const toggleFavorite = async (card: InfoCardData) => {
    const favoriteId = favoriteIds.get(card.id);
    if (!favoriteId) return;

    try {
      await api.removeFavorite(favoriteId);
      setFavorites(prev => prev.filter(f => f.id !== card.id));
      setFavoriteIds(prev => {
        const next = new Map(prev);
        next.delete(card.id);
        return next;
      });
      toast.success('已取消收藏');
    } catch (error) {
      console.error('Failed to remove favorite:', error);
      toast.error('取消收藏失败');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-background text-foreground flex items-center justify-center">
        <Loader2 size={40} className="animate-spin text-primary" />
      </div>
    );
  }

  if (!user) {
    return null; // Will redirect to login
  }

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
          <h2 className="text-2xl font-bold mb-1">{user.username}</h2>
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