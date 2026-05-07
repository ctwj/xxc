"use client";

import React, { useState, useRef, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Type, Bookmark, BookmarkCheck } from 'lucide-react';
import { cn } from '@/lib/utils';

interface LongTextReaderProps {
  title?: string;
  content: string;
  cardId: string;
}

type FontSize = 'small' | 'medium' | 'large';

const FONT_SIZE_MAP: Record<FontSize, string> = {
  small: 'text-sm md:text-base',
  medium: 'text-base md:text-lg',
  large: 'text-lg md:text-xl',
};

export const LongTextReader: React.FC<LongTextReaderProps> = ({
  title,
  content,
  cardId,
}) => {
  const [fontSize, setFontSize] = useState<FontSize>('medium');
  const [readProgress, setReadProgress] = useState(0);
  const [isBookmarked, setIsBookmarked] = useState(false);
  const contentRef = useRef<HTMLDivElement>(null);

  // 加载书签状态
  useEffect(() => {
    const bookmarks = JSON.parse(localStorage.getItem('reading-bookmarks') || '{}');
    setIsBookmarked(!!bookmarks[cardId]);
  }, [cardId]);

  // 监听滚动计算阅读进度
  const handleScroll = () => {
    if (!contentRef.current) return;

    const element = contentRef.current;
    const scrollTop = element.scrollTop;
    const scrollHeight = element.scrollHeight - element.clientHeight;
    const progress = scrollHeight > 0 ? (scrollTop / scrollHeight) * 100 : 0;

    setReadProgress(Math.min(progress, 100));

    // 自动保存阅读位置
    if (isBookmarked) {
      const bookmarks = JSON.parse(localStorage.getItem('reading-bookmarks') || '{}');
      bookmarks[cardId] = {
        scrollTop,
        timestamp: Date.now(),
      };
      localStorage.setItem('reading-bookmarks', JSON.stringify(bookmarks));
    }
  };

  // 恢复阅读位置
  useEffect(() => {
    if (!contentRef.current || !isBookmarked) return;

    const bookmarks = JSON.parse(localStorage.getItem('reading-bookmarks') || '{}');
    const bookmark = bookmarks[cardId];

    if (bookmark?.scrollTop) {
      setTimeout(() => {
        contentRef.current?.scrollTo({
          top: bookmark.scrollTop,
          behavior: 'smooth',
        });
      }, 100);
    }
  }, [cardId, isBookmarked]);

  const toggleBookmark = () => {
    const bookmarks = JSON.parse(localStorage.getItem('reading-bookmarks') || '{}');

    if (isBookmarked) {
      delete bookmarks[cardId];
      setIsBookmarked(false);
    } else {
      bookmarks[cardId] = {
        scrollTop: contentRef.current?.scrollTop || 0,
        timestamp: Date.now(),
      };
      setIsBookmarked(true);
    }

    localStorage.setItem('reading-bookmarks', JSON.stringify(bookmarks));
  };

  const cycleFontSize = () => {
    const sizes: FontSize[] = ['small', 'medium', 'large'];
    const currentIndex = sizes.indexOf(fontSize);
    const nextIndex = (currentIndex + 1) % sizes.length;
    setFontSize(sizes[nextIndex]);
  };

  return (
    <div className="relative flex flex-col h-full bg-gradient-to-b from-card to-background">
      {/* 阅读进度条 */}
      <div className="absolute top-0 left-0 right-0 h-1 bg-muted z-20">
        <motion.div
          className="h-full bg-primary"
          style={{ width: `${readProgress}%` }}
          transition={{ duration: 0.1 }}
        />
      </div>

      {/* 工具栏 */}
      <div className="absolute top-20 right-8 flex flex-col gap-2 z-10">
        {/* 字体大小调节 */}
        <motion.button
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
          onClick={cycleFontSize}
          className="p-3 bg-foreground/10 backdrop-blur-xl rounded-full text-foreground hover:bg-foreground/20 transition-all border border-foreground/10"
          title="调整字体大小"
        >
          <Type size={18} />
        </motion.button>

        {/* 书签 */}
        <motion.button
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
          onClick={toggleBookmark}
          className={cn(
            "p-3 backdrop-blur-xl rounded-full transition-all border",
            isBookmarked
              ? "bg-primary/20 text-primary border-primary/30"
              : "bg-foreground/10 text-foreground border-foreground/10 hover:bg-foreground/20"
          )}
          title={isBookmarked ? "取消书签" : "添加书签"}
        >
          {isBookmarked ? <BookmarkCheck size={18} /> : <Bookmark size={18} />}
        </motion.button>
      </div>

      {/* 内容区域 */}
      <div
        ref={contentRef}
        onScroll={handleScroll}
        className="flex-1 overflow-y-auto px-10 py-12 scroll-smooth"
      >
        {title && (
          <h2 className="text-3xl md:text-4xl font-bold mb-8 gradient-text leading-tight">
            {title}
          </h2>
        )}
        <div
          className={cn(
            "text-foreground/80 leading-loose space-y-4 transition-all duration-300",
            FONT_SIZE_MAP[fontSize]
          )}
        >
          {content.split('\n').map((paragraph, index) => (
            <p key={index} className="indent-8">
              {paragraph}
            </p>
          ))}
        </div>

        {/* 阅读完成提示 */}
        {readProgress >= 99 && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="mt-12 p-6 bg-primary/10 border border-primary/20 rounded-2xl text-center"
          >
            <p className="text-primary font-medium">
              ✨ 您已完成本文阅读
            </p>
          </motion.div>
        )}

        {/* 底部留白 */}
        <div className="h-20" />
      </div>
    </div>
  );
};