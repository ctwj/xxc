"use client";

import React, { useState, useRef, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Volume2, VolumeX, Clock, Share2, ChevronLeft, ChevronRight } from 'lucide-react';
import { InfoCardData } from '@/types';
import { cn } from '@/lib/utils';
import { formatRelativeTime } from '@/lib/utils';
import { ImagePreview } from './ImagePreview';
import { LongTextReader } from './LongTextReader';

interface InfoCardProps {
  data: InfoCardData;
  style?: React.CSSProperties;
  className?: string;
  isCurrent?: boolean;
  onShare?: () => void;
}

export const InfoCard: React.FC<InfoCardProps> = ({ data, style, className, isCurrent, onShare }) => {
  const [isMuted, setIsMuted] = useState(true);
  const [isExpanded, setIsExpanded] = useState(false);
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [showImagePreview, setShowImagePreview] = useState(false);
  const [previewImageIndex, setPreviewImageIndex] = useState(0);
  const videoRef = useRef<HTMLVideoElement>(null);

  // 离开当前卡片时自动暂停视频
  useEffect(() => {
    if (!isCurrent && videoRef.current) {
      videoRef.current.pause();
    } else if (isCurrent && videoRef.current && (data.type === 'video' || data.type === 'video-text' || data.type === 'images-video-text')) {
      videoRef.current.play().catch(e => console.log('Auto play blocked', e));
    }
  }, [isCurrent, data.type]);

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

  const renderContent = () => {
    // Handle html/markdown types as text
    const effectiveType = (data.type === 'html' || data.type === 'markdown') ? 'text' : data.type;

    switch (effectiveType) {
      case 'text':
        return (
          <div className="flex flex-col items-center justify-center h-full p-8 text-center bg-gradient-to-b from-card to-background">
            {data.title && (
              <h2 className="text-3xl md:text-4xl font-bold mb-6 gradient-text leading-tight">{data.title}</h2>
            )}
            <div className={cn(
              "text-lg md:text-xl text-foreground/80 leading-relaxed max-w-lg transition-all duration-300 overflow-y-auto max-h-[60vh] px-4",
              isExpanded ? "max-h-full" : ""
            )}>
              {data.content}
            </div>
            {data.content && data.content.length > 150 && !isExpanded && (
              <button
                onClick={() => setIsExpanded(true)}
                className="mt-4 text-primary font-medium hover:underline transition-all"
              >
                展开
              </button>
            )}
          </div>
        );

      case 'long-text':
        return (
          <LongTextReader
            title={data.title}
            content={data.content || ''}
            cardId={data.id}
          />
        );

      case 'image':
        return (
          <div className="relative w-full h-full overflow-hidden">
            <img
              src={data.mediaUrl}
              alt={data.title || '图片内容'}
              className="w-full h-full object-cover cursor-zoom-in"
              onClick={() => handleImageClick(0)}
              onError={(e) => {
                (e.target as HTMLImageElement).src = 'https://images.unsplash.com/photo-1557683316-973673baf926?q=80&w=1000&auto=format&fit=crop';
              }}
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent pointer-events-none" />
            <div className="absolute bottom-32 left-8 right-8 pointer-events-none">
              {data.title && <h2 className="text-3xl font-bold text-white mb-2">{data.title}</h2>}
              {data.content && <p className="text-white/80 line-clamp-2">{data.content}</p>}
            </div>
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
            <div className="relative h-[60%] w-full overflow-hidden">
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
            <div className="h-[40%] w-full p-8 flex flex-col justify-center overflow-y-auto">
              {data.title && <h2 className="text-2xl font-bold mb-4 gradient-text">{data.title}</h2>}
              {data.content && <p className="text-foreground/70 leading-relaxed">{data.content}</p>}
            </div>
          </div>
        );

      case 'images-video-text':
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
              <video
                ref={videoRef}
                src={data.videoUrl}
                className="w-full h-full object-cover"
                muted={isMuted}
                loop
                playsInline
                poster={data.coverUrl}
              />
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  setIsMuted(!isMuted);
                }}
                className="absolute top-2 right-2 p-2 bg-black/30 backdrop-blur-md rounded-full text-white/80 hover:bg-black/50 transition-colors z-10"
              >
                {isMuted ? <VolumeX size={16} /> : <Volume2 size={16} />}
              </button>
            </div>

            {/* 文字区域 */}
            <div className="h-[30%] w-full p-6 flex flex-col justify-center overflow-y-auto">
              {data.title && <h2 className="text-xl font-bold mb-3 gradient-text">{data.title}</h2>}
              {data.content && <p className="text-foreground/70 leading-relaxed text-sm line-clamp-3">{data.content}</p>}
            </div>
          </div>
        );

      case 'video':
      case 'video-text':
        return (
          <div className="relative w-full h-full overflow-hidden bg-black">
            <video
              ref={videoRef}
              src={data.mediaUrl}
              className={cn(
                "w-full h-full object-cover",
                data.type === 'video-text' ? "h-[65%]" : "h-full"
              )}
              muted={isMuted}
              loop
              playsInline
              poster={data.coverUrl}
            />
            <button
              onClick={(e) => {
                e.stopPropagation();
                setIsMuted(!isMuted);
              }}
              className="absolute top-8 right-8 p-3 bg-black/30 backdrop-blur-md rounded-full text-white/80 hover:bg-black/50 transition-colors z-10"
            >
              {isMuted ? <VolumeX size={20} /> : <Volume2 size={20} />}
            </button>

            {data.type === 'video' && (
              <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent pointer-events-none" />
            )}

            <div className={cn(
              "absolute left-0 right-0 p-8 flex flex-col justify-end",
              data.type === 'video-text'
                ? "bottom-0 h-[35%] bg-card"
                : "bottom-24 bg-gradient-to-t from-black/80 to-transparent"
            )}>
              {data.title && (
                <h2 className={cn(
                  "text-2xl font-bold mb-2",
                  data.type === 'video-text' ? "text-foreground" : "text-white"
                )}>{data.title}</h2>
              )}
              {data.content && (
                <p className={cn(
                  "line-clamp-3 text-sm md:text-base",
                  data.type === 'video-text' ? "text-foreground/70" : "text-white/70"
                )}>{data.content}</p>
              )}
            </div>
          </div>
        );

      case 'image-text':
        return (
          <div className="flex flex-col w-full h-full bg-card overflow-hidden">
            <div className="h-[60%] w-full overflow-hidden">
              <img
                src={data.mediaUrl}
                alt={data.title || '图片内容'}
                className="w-full h-full object-cover transition-transform duration-700 hover:scale-105 cursor-zoom-in"
                onClick={() => handleImageClick(0)}
              />
            </div>
            <div className="h-[40%] w-full p-8 flex flex-col justify-center">
              {data.title && <h2 className="text-2xl font-bold mb-4 gradient-text">{data.title}</h2>}
              {data.content && <p className="text-foreground/70 leading-relaxed line-clamp-4">{data.content}</p>}
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
        {/* 顶部时间信息 - 突出显示 */}
        <div className="absolute top-8 left-8 right-8 z-10 flex items-center justify-between">
          <div className="flex items-center gap-2 px-4 py-2 bg-foreground/40 backdrop-blur-xl rounded-full border border-foreground/10">
            <Clock size={16} className="text-primary" />
            <span className="text-sm font-bold text-background tracking-wide">
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
              className="p-3 bg-foreground/40 backdrop-blur-xl rounded-full text-background hover:bg-foreground/60 transition-all border border-foreground/10"
              title="分享"
            >
              <Share2 size={18} />
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
