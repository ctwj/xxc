"use client";

import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, ChevronLeft, ChevronRight, ZoomIn, ZoomOut } from 'lucide-react';
import { cn } from '@/lib/utils';

interface ImagePreviewProps {
  images: string[];
  initialIndex: number;
  isOpen: boolean;
  onClose: () => void;
}

export const ImagePreview: React.FC<ImagePreviewProps> = ({
  images,
  initialIndex,
  isOpen,
  onClose,
}) => {
  const [currentIndex, setCurrentIndex] = useState(initialIndex);
  const [scale, setScale] = useState(1);

  useEffect(() => {
    setCurrentIndex(initialIndex);
    setScale(1);
  }, [initialIndex, isOpen]);

  // 键盘导航
  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      } else if (e.key === 'ArrowLeft') {
        handlePrev();
      } else if (e.key === 'ArrowRight') {
        handleNext();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, currentIndex]);

  const handlePrev = () => {
    setCurrentIndex((prev) => (prev > 0 ? prev - 1 : images.length - 1));
    setScale(1);
  };

  const handleNext = () => {
    setCurrentIndex((prev) => (prev < images.length - 1 ? prev + 1 : 0));
    setScale(1);
  };

  const handleZoomIn = () => {
    setScale((prev) => Math.min(prev + 0.5, 3));
  };

  const handleZoomOut = () => {
    setScale((prev) => Math.max(prev - 0.5, 1));
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="fixed inset-0 z-[100] bg-black/95 backdrop-blur-sm flex items-center justify-center"
          onClick={onClose}
        >
          {/* 关闭按钮 */}
          <button
            onClick={onClose}
            className="absolute top-6 right-6 p-3 bg-white/10 backdrop-blur-md rounded-full text-white hover:bg-white/20 transition-all z-10"
          >
            <X size={24} />
          </button>

          {/* 缩放控制 */}
          <div className="absolute top-6 left-6 flex gap-2 z-10">
            <button
              onClick={(e) => {
                e.stopPropagation();
                handleZoomOut();
              }}
              disabled={scale <= 1}
              className={cn(
                "p-3 bg-white/10 backdrop-blur-md rounded-full text-white transition-all",
                scale <= 1 ? "opacity-50 cursor-not-allowed" : "hover:bg-white/20"
              )}
            >
              <ZoomOut size={20} />
            </button>
            <button
              onClick={(e) => {
                e.stopPropagation();
                handleZoomIn();
              }}
              disabled={scale >= 3}
              className={cn(
                "p-3 bg-white/10 backdrop-blur-md rounded-full text-white transition-all",
                scale >= 3 ? "opacity-50 cursor-not-allowed" : "hover:bg-white/20"
              )}
            >
              <ZoomIn size={20} />
            </button>
            <div className="px-4 py-3 bg-white/10 backdrop-blur-md rounded-full text-white text-sm font-medium">
              {Math.round(scale * 100)}%
            </div>
          </div>

          {/* 图片计数 */}
          {images.length > 1 && (
            <div className="absolute bottom-6 left-1/2 -translate-x-1/2 px-4 py-2 bg-white/10 backdrop-blur-md rounded-full text-white text-sm font-medium z-10">
              {currentIndex + 1} / {images.length}
            </div>
          )}

          {/* 左右导航按钮 */}
          {images.length > 1 && (
            <>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  handlePrev();
                }}
                className="absolute left-6 top-1/2 -translate-y-1/2 p-4 bg-white/10 backdrop-blur-md rounded-full text-white hover:bg-white/20 transition-all z-10"
              >
                <ChevronLeft size={28} />
              </button>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  handleNext();
                }}
                className="absolute right-6 top-1/2 -translate-y-1/2 p-4 bg-white/10 backdrop-blur-md rounded-full text-white hover:bg-white/20 transition-all z-10"
              >
                <ChevronRight size={28} />
              </button>
            </>
          )}

          {/* 图片展示 */}
          <motion.div
            key={currentIndex}
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.9 }}
            transition={{ duration: 0.2 }}
            className="relative max-w-[90vw] max-h-[90vh] overflow-auto"
            onClick={(e) => e.stopPropagation()}
          >
            <motion.img
              src={images[currentIndex]}
              alt={`预览图片 ${currentIndex + 1}`}
              className="w-auto h-auto max-w-full max-h-[90vh] object-contain transition-transform duration-300"
              style={{ transform: `scale(${scale})` }}
              draggable={false}
            />
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};