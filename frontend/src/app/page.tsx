"use client";

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { User, LogOut, LogIn, ChevronLeft, ChevronRight, Sun, Moon } from 'lucide-react';
import { CardStack } from '@/components/home/CardStack';
import { ActionBar } from '@/components/home/ActionBar';
import { InfoCardData, Article, ContentType } from '@/types';
import { motion, AnimatePresence } from 'framer-motion';
import { toast } from 'sonner';
import { useTheme } from '@/contexts/ThemeContext';
import { api } from '@/lib/api';

// Convert Article to InfoCardData
function articleToCardData(article: Article): InfoCardData {
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

  // Determine content type based on media fields
  // Note: article.type from API is content format (html/markdown), not display type
  let type: ContentType = 'text';

  if (article.videoUrl) {
    // Has video
    type = article.description ? 'video-text' : 'video';
  } else if (mediaUrls.length > 1) {
    // Multiple images
    type = article.description ? 'images-text' : 'images';
  } else if (article.thumbnail || mediaUrls.length === 1) {
    // Single image
    type = article.description ? 'image-text' : 'image';
  } else if (article.description && article.description.length > 500) {
    // Long text without media
    type = 'long-text';
  }

  return {
    id: String(article.id),
    type,
    title: article.title,
    content: article.description,
    mediaUrl: article.thumbnail,
    mediaUrls,
    videoUrl: article.videoUrl,
    coverUrl: article.coverUrl,
    tag: article.tag || '未分类',
    publishTime: new Date(article.createTime * 1000).toISOString(),
    isNew: Date.now() - article.createTime * 1000 < 24 * 60 * 60 * 1000, // Within 24 hours
  };
}

export default function HomePage() {
  const router = useRouter();
  const { theme, toggleTheme } = useTheme();
  const [user, setUser] = useState<{ id: number; username: string; role: string } | null>(null);
  const [cards, setCards] = useState<InfoCardData[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [favorites, setFavorites] = useState<Set<string>>(new Set());
  const [likeStatus, setLikeStatus] = useState<Map<string, { type: number; likes: number; dislikes: number }>>(new Map());
  const [loading, setLoading] = useState(true);

  // 边界状态
  const isAtNewerBoundary = currentIndex === -1;
  const isAtOlderBoundary = currentIndex === cards.length;

  // Load initial data
  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        // Load articles
        const response = await api.getArticles({ page: 1, pageSize: 50 });
        const cardData = response.data.map(articleToCardData);
        setCards(cardData);

        // Check if user is logged in
        try {
          const userData = await api.getCurrentUser();
          setUser(userData);

          // Load favorites for logged-in user
          const favResponse = await api.getFavorites();
          const favIds = new Set(favResponse.data.map(f => String(f.articleId)));
          setFavorites(favIds);
        } catch {
          // Not logged in, continue without user data
          setUser(null);
        }
      } catch (error) {
        console.error('Failed to load data:', error);
        toast.error('加载数据失败');
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  // Load like status for current card
  useEffect(() => {
    if (!user || cards.length === 0) return;

    const loadLikeStatus = async () => {
      const currentCard = cards[currentIndex];
      if (!currentCard || currentIndex < 0) return;

      try {
        const status = await api.getLikeStatus(Number(currentCard.id));
        setLikeStatus(prev => new Map(prev).set(currentCard.id, status));
      } catch (error) {
        console.error('Failed to load like status:', error);
      }
    };

    loadLikeStatus();
  }, [user, currentIndex, cards]);

  const handleRefresh = async () => {
    setLoading(true);
    try {
      const response = await api.getArticles({ page: 1, pageSize: 50 });
      const cardData = response.data.map(articleToCardData);
      setCards(cardData);
      setCurrentIndex(0);
      toast.success('刷新成功');
    } catch (error) {
      toast.error('刷新失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSignOut = async () => {
    try {
      await api.logout();
      setUser(null);
      setFavorites(new Set());
      setLikeStatus(new Map());
      toast.success('已退出登录');
    } catch (error) {
      toast.error('退出失败');
    }
  };

  // Navigation
  const handleViewOlder = () => {
    if (isAtNewerBoundary) {
      setCurrentIndex(0);
      return;
    }

    if (currentIndex < cards.length - 1) {
      setCurrentIndex(prev => prev + 1);
    } else if (currentIndex === cards.length - 1) {
      setCurrentIndex(cards.length);
    }
  };

  const handleViewNewer = () => {
    if (isAtOlderBoundary) {
      setCurrentIndex(cards.length - 1);
      return;
    }

    if (currentIndex > 0) {
      setCurrentIndex(prev => prev - 1);
    } else if (currentIndex === 0) {
      setCurrentIndex(-1);
    }
  };

  const handleLike = async () => {
    const currentCard = cards[currentIndex];
    if (!user || !currentCard) {
      handleViewNewer();
      return;
    }

    try {
      await api.setLike(Number(currentCard.id), 1);
      // Update local status
      const status = likeStatus.get(currentCard.id) || { type: 0, likes: 0, dislikes: 0 };
      setLikeStatus(prev => new Map(prev).set(currentCard.id, {
        type: status.type === 1 ? 0 : 1,
        likes: status.type === 1 ? status.likes - 1 : status.likes + 1,
        dislikes: status.dislikes,
      }));
    } catch (error) {
      console.error('Failed to like:', error);
    }

    handleViewNewer();
  };

  const handleDislike = async () => {
    const currentCard = cards[currentIndex];
    if (!user || !currentCard) {
      handleViewOlder();
      return;
    }

    try {
      await api.setLike(Number(currentCard.id), 2);
      // Update local status
      const status = likeStatus.get(currentCard.id) || { type: 0, likes: 0, dislikes: 0 };
      setLikeStatus(prev => new Map(prev).set(currentCard.id, {
        type: status.type === 2 ? 0 : 2,
        likes: status.likes,
        dislikes: status.type === 2 ? status.dislikes - 1 : status.dislikes + 1,
      }));
    } catch (error) {
      console.error('Failed to dislike:', error);
    }

    handleViewOlder();
  };

  const handleFavorite = async () => {
    if (!user) {
      toast.error('请先登录');
      router.push('/login');
      return;
    }

    const currentCard = cards[currentIndex];
    if (!currentCard) return;

    try {
      if (favorites.has(currentCard.id)) {
        // Find favorite ID and remove
        const favResponse = await api.getFavorites();
        const fav = favResponse.data.find(f => String(f.articleId) === currentCard.id);
        if (fav) {
          await api.removeFavorite(fav.id);
          setFavorites(prev => {
            const next = new Set(prev);
            next.delete(currentCard.id);
            return next;
          });
          toast.success('已取消收藏');
        }
      } else {
        await api.addFavorite(Number(currentCard.id));
        setFavorites(prev => new Set(prev).add(currentCard.id));
        toast.success('已添加到收藏');
      }
    } catch (error) {
      toast.error('操作失败');
    }
  };

  const handleShareCard = (card: InfoCardData) => {
    const shareText = `${card.title || '精彩内容'}\n\n${card.content || ''}\n\n来自每日信息差`;

    if (navigator.share) {
      navigator.share({
        title: card.title || '每日信息差',
        text: shareText,
        url: window.location.href,
      }).catch(() => {});
    } else {
      navigator.clipboard.writeText(shareText).then(() => {
        toast.success('内容已复制到剪贴板');
      }).catch(() => {
        toast.error('分享失败');
      });
    }
  };

  const currentCard = cards[currentIndex];

  if (loading) {
    return (
      <div className="relative h-screen w-screen bg-background overflow-hidden flex items-center justify-center">
        <div className="text-xl font-bold gradient-text">加载中...</div>
      </div>
    );
  }

  return (
    <div className="relative h-screen w-screen bg-background overflow-hidden flex flex-col">
      {/* 顶部状态栏 */}
      <header className="fixed top-0 left-0 right-0 z-50 px-8 py-10 flex items-center justify-between pointer-events-none">
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          className="text-xl font-bold tracking-tight text-foreground pointer-events-auto"
        >
          每日信息差
        </motion.div>
        <div className="flex items-center gap-3 pointer-events-auto">
          <motion.button
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            onClick={toggleTheme}
            className="p-3 bg-foreground/10 backdrop-blur-md rounded-full text-foreground hover:bg-foreground/20 transition-all border border-foreground/10"
          >
            {theme === 'dark' ? <Sun size={22} /> : <Moon size={22} />}
          </motion.button>

          {user ? (
            <>
              <motion.button
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                onClick={() => router.push('/favorites')}
                className="p-3 bg-foreground/10 backdrop-blur-md rounded-full text-foreground hover:bg-foreground/20 transition-all border border-foreground/10"
              >
                <User size={22} />
              </motion.button>
              <motion.button
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                onClick={handleSignOut}
                className="p-3 bg-foreground/10 backdrop-blur-md rounded-full text-foreground hover:bg-destructive/80 transition-all border border-foreground/10"
              >
                <LogOut size={22} />
              </motion.button>
            </>
          ) : (
            <motion.button
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              onClick={() => router.push('/login')}
              className="px-5 py-2.5 bg-primary backdrop-blur-md rounded-full text-primary-foreground hover:bg-primary/90 transition-all border border-primary/20 font-semibold text-sm flex items-center gap-2"
            >
              <LogIn size={18} />
              登录
            </motion.button>
          )}
        </div>
      </header>

      {/* 核心卡片堆叠 */}
      <main className="flex-1 flex items-center justify-center relative">
        {/* 没有更新的消息提示 */}
        {isAtNewerBoundary && (
          <motion.div
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            className="absolute inset-0 flex items-center justify-center z-30"
          >
            <div className="glass-card rounded-[20px] p-12 max-w-md mx-4 text-center">
              <h2 className="text-2xl font-bold mb-4 text-foreground">当前没有更新的消息</h2>
              <p className="text-foreground/60 mb-8 leading-relaxed">
                您已经浏览到最新的内容了。请等待消息更新，或者浏览历史消息。
              </p>
              <button
                onClick={handleViewOlder}
                className="px-8 py-3 bg-primary text-primary-foreground rounded-full font-semibold hover:bg-primary/90 transition-all"
              >
                浏览历史消息
              </button>
            </div>
          </motion.div>
        )}

        {/* 需要登录提示 */}
        {isAtOlderBoundary && (
          <motion.div
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            className="absolute inset-0 flex items-center justify-center z-30"
          >
            <div className="glass-card rounded-[20px] p-12 max-w-md mx-4 text-center">
              <h2 className="text-2xl font-bold mb-4 text-foreground">登录后可以查看更多信息</h2>
              <p className="text-foreground/60 mb-8 leading-relaxed">
                您已浏览完所有内容。登录后可以享受收藏、点赞等更多功能。
              </p>
              <div className="flex gap-4 justify-center">
                <button
                  onClick={() => router.push('/login')}
                  className="px-8 py-3 bg-primary text-primary-foreground rounded-full font-semibold hover:bg-primary/90 transition-all"
                >
                  立即登录
                </button>
                <button
                  onClick={handleViewNewer}
                  className="px-8 py-3 bg-foreground/10 text-foreground rounded-full font-semibold hover:bg-foreground/20 transition-all"
                >
                  返回最后一条
                </button>
              </div>
            </div>
          </motion.div>
        )}

        {/* 正常卡片显示 */}
        {!isAtNewerBoundary && !isAtOlderBoundary && cards.length > 0 && (
          <CardStack
            cards={cards}
            currentIndex={currentIndex}
            onSwipe={() => {}}
            onFavorite={() => {}}
            isFavorited={(id) => favorites.has(id)}
            onRefresh={handleRefresh}
            onShare={handleShareCard}
          />
        )}

        {/* 左右箭头按钮 */}
        {cards.length > 0 && !isAtOlderBoundary && (
          <motion.button
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 0.6, x: 0 }}
            whileHover={{ opacity: 1, scale: 1.1 }}
            onClick={handleDislike}
            className="hidden md:flex absolute left-8 top-1/2 -translate-y-1/2 z-40 w-14 h-14 items-center justify-center bg-foreground/10 backdrop-blur-md rounded-full text-foreground hover:bg-foreground/20 transition-all border border-foreground/10"
          >
            <ChevronLeft size={28} />
          </motion.button>
        )}

        {cards.length > 0 && !isAtNewerBoundary && (
          <motion.button
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 0.6, x: 0 }}
            whileHover={{ opacity: 1, scale: 1.1 }}
            onClick={handleLike}
            className="hidden md:flex absolute right-8 top-1/2 -translate-y-1/2 z-40 w-14 h-14 items-center justify-center bg-foreground/10 backdrop-blur-md rounded-full text-foreground hover:bg-foreground/20 transition-all border border-foreground/10"
          >
            <ChevronRight size={28} />
          </motion.button>
        )}
      </main>

      {/* 底部操作栏 */}
      <AnimatePresence>
        {cards.length > 0 && !isAtNewerBoundary && !isAtOlderBoundary && (
          <motion.div
            initial={{ y: 100, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            exit={{ y: 100, opacity: 0 }}
          >
            <ActionBar
              onLike={handleLike}
              onDislike={handleDislike}
              onFavorite={handleFavorite}
              isFavorited={currentCard ? favorites.has(currentCard.id) : false}
            />
          </motion.div>
        )}
      </AnimatePresence>

      {/* 背景动态色块 */}
      <div className="absolute inset-0 z-[-1] overflow-hidden pointer-events-none opacity-30">
        <div className="absolute top-[-10%] left-[-10%] w-[50%] h-[50%] bg-primary/20 blur-[120px] rounded-full animate-pulse" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[50%] h-[50%] bg-accent/10 blur-[120px] rounded-full animate-pulse" />
      </div>
    </div>
  );
}