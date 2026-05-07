"use client";

import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Clock, Share2, ChevronLeft, ChevronRight } from 'lucide-react';
import dynamic from 'next/dynamic';
import { InfoCardData } from '@/types';
import { cn } from '@/lib/utils';
import { formatRelativeTime } from '@/lib/utils';
import { ImagePreview } from './ImagePreview';
import { LongTextReader } from './LongTextReader';

// 动态导入 Plyr，避免 SSR 问题
const Plyr = dynamic(() => import('plyr-react').then(mod => ({ default: mod.Plyr })), {
  ssr: false,
  loading: () => <div className="w-full h-full bg-black flex items-center justify-center text-white">加载播放器...</div>
});

interface InfoCardProps {
  data: InfoCardData;
  style?: React.CSSProperties;
  className?: string;
  isCurrent?: boolean;
  onShare?: () => void;
}

// 根据文字长度计算字体大小
const getFontSizeClass = (textLength: number): string => {
  if (textLength <= 20) return 'text-2xl md:text-3xl';
  if (textLength <= 50) return 'text-xl md:text-2xl';
  if (textLength <= 100) return 'text-lg md:text-xl';
  if (textLength <= 200) return 'text-base md:text-lg';
  return 'text-sm md:text-base';
};

// 检测图片方向（需要异步加载图片获取尺寸）
const useImageOrientation = (url: string | undefined): 'portrait' | 'landscape' | 'unknown' => {
  const [orientation, setOrientation] = React.useState<'portrait' | 'landscape' | 'unknown'>('unknown');

  useEffect(() => {
    if (!url) {
      setOrientation('unknown');
      return;
    }

    const img = new Image();
    img.onload = () => {
      const isPortrait = img.height > img.width;
      setOrientation(isPortrait ? 'portrait' : 'landscape');
    };
    img.onerror = () => setOrientation('unknown');
    img.src = url;
  }, [url]);

  return orientation;
};

export const InfoCard: React.FC<InfoCardProps> = ({ data, style, className, isCurrent, onShare }) => {
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [showImagePreview, setShowImagePreview] = useState(false);
  const [previewImageIndex, setPreviewImageIndex] = useState(0);

  // 获取单图的方向
  const singleImageOrientation = useImageOrientation(data.mediaUrl);

  // 重置图片索引
  useEffect(() => {
    setCurrentImageIndex(0);
  }, [data.id]);

  const handlePrevImage = (e: React.MouseEvent) => {
    e.stopPropagation();
    setCurrentImageIndex(prev => prev > 0 ? prev - 1 : (data.mediaUrls?.length || 1) - 1);
  };

  const handleNextImage = (e: React.MouseEvent) => {
    e.stopPropagation();
    setCurrentImageIndex(prev => prev < (data.mediaUrls?.length || 1) - 1 ? prev + 1 : 0);
  };

  const handleImageClick = (index: number) => {
    setPreviewImageIndex(index);
    setShowImagePreview(true);
  };

  const getImageUrls = (): string[] => {
    if (data.mediaUrls && data.mediaUrls.length > 0) {
      return data.mediaUrls;
    }
    if (data.mediaUrl) {
      return [data.mediaUrl];
    }
    return [];
  };

  // 获取文字内容长度
  const contentLength = data.content?.length || 0;
  const fontSizeClass = getFontSizeClass(contentLength);

  const renderContent = () => {
    // Handle html/markdown types as text
    const effectiveType = (data.type === 'html' || data.type === 'markdown') ? 'text' : data.type;

    switch (effectiveType) {
      case 'text':
        return (
          <div className="flex flex-col items-center justify-center h-full p-8 text-center bg-gradient-to-b from-card to-background">
            <div className={cn(
              `${fontSizeClass} text-foreground/90 leading-relaxed max-w-lg overflow-y-auto px-4`,
              contentLength > 200 ? "max-h-[60vh]" : ""
            )}>
              {data.content}
            </div>
          </div>
        );

      case 'long-text':
        return (
          <LongTextReader
            title=""
            content={data.content || ''}
            cardId={data.id}
          />
        );

      case 'image':
        // 竖图：图片占满，文字在下方有底色
        if (singleImageOrientation === 'portrait') {
          return (
            <div className="relative w-full h-full overflow-hidden">
              <img
                src={data.mediaUrl}
                alt="图片内容"
                className="w-full h-full object-cover cursor-zoom-in"
                onClick={() => handleImageClick(0)}
                onError={(e) => {
                  (e.target as HTMLImageElement).src = 'https://images.unsplash.com/photo-1557683316-973673baf926?q=80&w=1000&auto=format&fit=crop';
                }}
              />
              {data.content && (
                <div className="absolute bottom-0 left-0 right-0 p-6 bg-gradient-to-t from-black/90 via-black/70 to-transparent">
                  <p className={cn(
                    `${fontSizeClass} text-white leading-relaxed`,
                    contentLength > 100 ? "line-clamp-3" : ""
                  )}>
                    {data.content}
                  </p>
                </div>
              )}
            </div>
          );
        }
        // 横图或未知：保持当前布局（图片上方，文字下方）
        return (
          <div className="relative w-full h-full overflow-hidden">
            <img
              src={data.mediaUrl}
              alt="图片内容"
              className="w-full h-full object-cover cursor-zoom-in"
              onClick={() => handleImageClick(0)}
              onError={(e) => {
                (e.target as HTMLImageElement).src = 'https://images.unsplash.com/photo-1557683316-973673baf926?q=80&w=1000&auto=format&fit=crop';
              }}
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent pointer-events-none" />
            {data.content && (
              <div className="absolute bottom-32 left-8 right-8 pointer-events-none">
                <p className={cn(
                  `${fontSizeClass} text-white leading-relaxed`,
                  contentLength > 100 ? "line-clamp-2" : ""
                )}>
                  {data.content}
                </p>
              </div>
            )}
          </div>
        );

      case 'images':
        return (
          <div className="relative w-full h-full overflow-hidden bg-black">
            {data.mediaUrls && data.mediaUrls.length > 0 && (
              <>
                <img
                  src={data.mediaUrls[currentImageIndex]}
                  alt={`图片 ${currentImageIndex + 1}`}
                  className="w-full h-full object-cover transition-opacity duration-300 cursor-zoom-in"
                  onClick={() => handleImageClick(currentImageIndex)}
                />

                {/* 图片导航按钮 */}
                {data.mediaUrls.length > 1 && (
                  <>
                    <button
                      onClick={handlePrevImage}
                      className="absolute left-4 top-1/2 -translate-y-1/2 p-2 bg-black/50 backdrop-blur-md rounded-full text-white hover:bg-black/70 transition-all z-10"
                    >
                      <ChevronLeft size={24} />
                    </button>
                    <button
                      onClick={handleNextImage}
                      className="absolute right-4 top-1/2 -translate-y-1/2 p-2 bg-black/50 backdrop-blur-md rounded-full text-white hover:bg-black/70 transition-all z-10"
                    >
                      <ChevronRight size={24} />
                    </button>

                    {/* 图片指示器 */}
                    <div className="absolute bottom-8 left-1/2 -translate-x-1/2 flex gap-2 z-10">
                      {data.mediaUrls.map((_, index) => (
                        <div
                          key={index}
                          className={cn(
                            "w-2 h-2 rounded-full transition-all",
                            index === currentImageIndex ? "bg-white w-6" : "bg-white/50"
                          )}
                        />
                      ))}
                    </div>
                  </>
                )}
              </>
            )}
          </div>
        );

      case 'images-text':
        return (
          <div className="flex flex-col w-full h-full bg-card overflow-hidden">
            <div className="relative h-[55%] w-full overflow-hidden">
              {data.mediaUrls && data.mediaUrls.length > 0 && (
                <>
                  <img
                    src={data.mediaUrls[currentImageIndex]}
                    alt={`图片 ${currentImageIndex + 1}`}
                    className="w-full h-full object-cover transition-opacity duration-300 cursor-zoom-in"
                    onClick={() => handleImageClick(currentImageIndex)}
                  />

                  {/* 图片导航 */}
                  {data.mediaUrls.length > 1 && (
                    <>
                      <button
                        onClick={handlePrevImage}
                        className="absolute left-4 top-1/2 -translate-y-1/2 p-2 bg-black/50 backdrop-blur-md rounded-full text-white hover:bg-black/70 transition-all z-10"
                      >
                        <ChevronLeft size={20} />
                      </button>
                      <button
                        onClick={handleNextImage}
                        className="absolute right-4 top-1/2 -translate-y-1/2 p-2 bg-black/50 backdrop-blur-md rounded-full text-white hover:bg-black/70 transition-all z-10"
                      >
                        <ChevronRight size={20} />
                      </button>

                      <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2 z-10">
                        {data.mediaUrls.map((_, index) => (
                          <div
                            key={index}
                            className={cn(
                              "w-1.5 h-1.5 rounded-full transition-all",
                              index === currentImageIndex ? "bg-white w-4" : "bg-white/50"
                            )}
                          />
                        ))}
                      </div>
                    </>
                  )}
                </>
              )}
            </div>
            {/* 文字区域 - 根据内容长度调整样式 */}
            <div className="h-[45%] w-full flex flex-col justify-center overflow-y-auto px-8 py-6 bg-gradient-to-b from-card to-background">
              {data.content && (
                <p className={cn(
                  `${fontSizeClass} text-foreground/80 leading-relaxed text-center`
                )}>
                  {data.content}
                </p>
              )}
            </div>
          </div>
        );

      case 'images-video-text':
        const imagesVideoOptions = {
          controls: ['play-large', 'play', 'progress', 'current-time', 'duration', 'mute', 'volume', 'fullscreen'],
          autoplay: false,
          muted: true,
          loop: { active: true },
        };
        const imagesVideoSrc = {
          type: 'video' as const,
          sources: [{ src: data.videoUrl || '', type: 'video/mp4' }],
          poster: data.coverUrl || '',
        };
        return (
          <div className="flex flex-col w-full h-full bg-card overflow-hidden">
            {/* 图片轮播区域 */}
            <div className="relative h-[35%] w-full overflow-hidden">
              {data.mediaUrls && data.mediaUrls.length > 0 && (
                <>
                  <img
                    src={data.mediaUrls[currentImageIndex]}
                    alt={`图片 ${currentImageIndex + 1}`}
                    className="w-full h-full object-cover cursor-zoom-in"
                    onClick={() => handleImageClick(currentImageIndex)}
                  />

                  {data.mediaUrls.length > 1 && (
                    <>
                      <button
                        onClick={handlePrevImage}
                        className="absolute left-2 top-1/2 -translate-y-1/2 p-1.5 bg-black/50 backdrop-blur-md rounded-full text-white hover:bg-black/70 transition-all z-10"
                      >
                        <ChevronLeft size={18} />
                      </button>
                      <button
                        onClick={handleNextImage}
                        className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 bg-black/50 backdrop-blur-md rounded-full text-white hover:bg-black/70 transition-all z-10"
                      >
                        <ChevronRight size={18} />
                      </button>

                      <div className="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1.5 z-10">
                        {data.mediaUrls.map((_, index) => (
                          <div
                            key={index}
                            className={cn(
                              "w-1 h-1 rounded-full transition-all",
                              index === currentImageIndex ? "bg-white w-3" : "bg-white/50"
                            )}
                          />
                        ))}
                      </div>
                    </>
                  )}
                </>
              )}
            </div>

            {/* 视频区域 */}
            <div className="relative h-[35%] w-full overflow-hidden bg-black">
              <Plyr
                options={imagesVideoOptions}
                source={imagesVideoSrc}
              />
            </div>

            {/* 文字区域 */}
            <div className="h-[30%] w-full p-6 flex flex-col justify-center overflow-y-auto">
              {data.content && (
                <p className={cn(
                  `${fontSizeClass} text-foreground/70 leading-relaxed`,
                  contentLength > 100 ? "line-clamp-3" : ""
                )}>
                  {data.content}
                </p>
              )}
            </div>
          </div>
        );

      case 'video':
      case 'video-text':
        // 视频源：优先使用 videoUrl，否则使用 mediaUrl
        const videoSource = data.videoUrl || data.mediaUrl || '';
        const videoOptions = {
          controls: ['play-large', 'play', 'progress', 'current-time', 'duration', 'mute', 'volume', 'settings', 'fullscreen'],
          settings: ['quality', 'speed'],
          autoplay: isCurrent,
          muted: true,
          loop: { active: true },
        };
        const videoSrc = {
          type: 'video' as const,
          sources: [
            {
              src: videoSource as string,
              type: 'video/mp4',
            },
          ],
          poster: (data.coverUrl || data.mediaUrl || '') as string,
        };

        // video-text 类型：上下分栏布局，文字紧贴视频
        if (data.type === 'video-text') {
          return (
            <div className="flex flex-col w-full h-full bg-card overflow-hidden">
              <div className="flex-shrink-0 w-full" style={{ height: 'auto', maxHeight: '60%' }}>
                <Plyr options={videoOptions} source={videoSrc} />
              </div>
              <div className="flex-1 w-full p-6 flex flex-col justify-center overflow-y-auto">
                {data.content && (
                  <p className={cn(
                    `${fontSizeClass} text-foreground/70 leading-relaxed`,
                    contentLength > 100 ? "line-clamp-4" : ""
                  )}>
                    {data.content}
                  </p>
                )}
              </div>
            </div>
          );
        }

        // 纯视频类型
        return (
          <div className="relative w-full h-full overflow-hidden bg-black">
            <div className="w-full h-full">
              <Plyr options={videoOptions} source={videoSrc} />
            </div>
            <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent pointer-events-none" />
            <div className="absolute left-0 right-0 bottom-24 p-8 flex flex-col justify-end pointer-events-none">
              {data.content && (
                <p className={cn(
                  `${fontSizeClass} text-white leading-relaxed`,
                  contentLength > 100 ? "line-clamp-3" : ""
                )}>
                  {data.content}
                </p>
              )}
            </div>
          </div>
        );

      case 'image-text':
        // 竖图：图片占满，文字在下方有底色
        if (singleImageOrientation === 'portrait') {
          return (
            <div className="relative w-full h-full overflow-hidden">
              <img
                src={data.mediaUrl}
                alt="图片内容"
                className="w-full h-full object-cover transition-transform duration-700 hover:scale-[1.02] cursor-zoom-in"
                onClick={() => handleImageClick(0)}
              />
              {data.content && (
                <div className="absolute bottom-0 left-0 right-0 p-6 bg-gradient-to-t from-black/90 via-black/70 to-transparent">
                  <p className={cn(
                    `${fontSizeClass} text-white leading-relaxed`,
                    contentLength > 100 ? "line-clamp-3" : ""
                  )}>
                    {data.content}
                  </p>
                </div>
              )}
            </div>
          );
        }
        // 横图：图片上方，文字下方
        return (
          <div className="flex flex-col w-full h-full bg-card overflow-hidden">
            <div className="h-[60%] w-full overflow-hidden">
              <img
                src={data.mediaUrl}
                alt="图片内容"
                className="w-full h-full object-cover transition-transform duration-700 hover:scale-[1.02] cursor-zoom-in"
                onClick={() => handleImageClick(0)}
              />
            </div>
            <div className="h-[40%] w-full flex flex-col justify-center overflow-y-auto px-8 py-6 bg-gradient-to-b from-card to-background">
              {data.content && (
                <p className={cn(
                  `${fontSizeClass} text-foreground/80 leading-relaxed text-center`
                )}>
                  {data.content}
                </p>
              )}
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <>
      <div
        className={cn(
          "relative w-full h-full rounded-[20px] overflow-hidden bg-card shadow-2xl transition-transform duration-500",
          className
        )}
        style={style}
      >
        {/* 顶部时间信息 - 简洁优雅 */}
        <div className="absolute top-6 left-6 right-6 z-10 flex items-center justify-between">
          <div className="flex items-center gap-2 px-3 py-1.5 bg-black/30 backdrop-blur-xl rounded-full border border-white/10">
            <Clock size={14} className="text-white/80" />
            <span className="text-xs font-medium text-white/90 tracking-wide">
              {formatRelativeTime(data.publishTime)}
            </span>
          </div>

          {/* 分享按钮 - 右上角 */}
          {onShare && (
            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={(e) => {
                e.stopPropagation();
                onShare();
              }}
              className="p-2.5 bg-black/30 backdrop-blur-xl rounded-full text-white/80 hover:bg-black/50 hover:text-white transition-all border border-white/10"
              title="分享"
            >
              <Share2 size={16} />
            </motion.button>
          )}
        </div>

        {/* 核心内容 */}
        {renderContent()}

        {/* 底部渐变增强层次感 */}
        {!['video-text', 'image-text', 'images-text', 'images-video-text', 'long-text'].includes(data.type) && (
           <div className="absolute bottom-0 left-0 right-0 h-32 bg-gradient-to-t from-black/40 to-transparent pointer-events-none" />
        )}
      </div>

      {/* 图片预览组件 */}
      <ImagePreview
        images={getImageUrls()}
        initialIndex={previewImageIndex}
        isOpen={showImagePreview}
        onClose={() => setShowImagePreview(false)}
      />
    </>
  );
};