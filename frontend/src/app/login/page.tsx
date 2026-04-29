"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { motion } from 'framer-motion';
import { Mail, Lock, ArrowLeft, Loader2 } from 'lucide-react';
import { toast } from 'sonner';

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!email.trim() || !password.trim()) {
      toast.error('请输入邮箱和密码');
      return;
    }

    // 验证邮箱格式
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      toast.error('请输入有效的邮箱地址');
      return;
    }

    setLoading(true);
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username: email, password }),
        credentials: 'include',
      });

      const data = await response.json();

      if (response.ok && data.success) {
        toast.success('登录成功');
        router.push('/');
        router.refresh();
      } else {
        toast.error('登录失败：' + (data.error || '未知错误'));
      }
    } catch (err) {
      toast.error('网络错误，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative min-h-screen w-screen bg-background overflow-hidden flex items-center justify-center">
      {/* 背景动态色块 */}
      <div className="absolute inset-0 z-0 overflow-hidden pointer-events-none opacity-30">
        <div className="absolute top-[-10%] left-[-10%] w-[50%] h-[50%] bg-primary/20 blur-[120px] rounded-full animate-pulse" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[50%] h-[50%] bg-accent/10 blur-[120px] rounded-full animate-pulse delay-1000" />
      </div>

      {/* 返回按钮 */}
      <button
        onClick={() => router.push('/')}
        className="fixed top-8 left-8 z-50 p-3 bg-white/10 backdrop-blur-md rounded-full text-white hover:bg-white/20 transition-all border border-white/10 active:scale-90"
      >
        <ArrowLeft size={22} />
      </button>

      {/* 登录表单卡片 */}
      <motion.div
        initial={{ scale: 0.9, opacity: 0, y: 20 }}
        animate={{ scale: 1, opacity: 1, y: 0 }}
        transition={{ type: 'spring', damping: 20, stiffness: 300 }}
        className="relative z-10 w-full max-w-md mx-4"
      >
        <div className="glass-card rounded-[20px] p-10 shadow-2xl">
          {/* 标题 */}
          <div className="text-center mb-10">
            <h1 className="text-3xl font-bold mb-2 gradient-text">欢迎回来</h1>
            <p className="text-foreground/60 text-sm">登录以继续使用每日信息差</p>
          </div>

          {/* 表单 */}
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="space-y-2">
              <label htmlFor="email" className="text-foreground/80 text-sm font-medium block">
                邮箱地址
              </label>
              <div className="relative">
                <Mail size={18} className="absolute left-4 top-1/2 -translate-y-1/2 text-foreground/40" />
                <input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="your@email.com"
                  className="w-full pl-12 h-12 bg-background/50 border border-white/10 rounded-full text-foreground placeholder:text-foreground/30 focus:border-primary transition-all px-4"
                  disabled={loading}
                />
              </div>
            </div>

            <div className="space-y-2">
              <label htmlFor="password" className="text-foreground/80 text-sm font-medium block">
                密码
              </label>
              <div className="relative">
                <Lock size={18} className="absolute left-4 top-1/2 -translate-y-1/2 text-foreground/40" />
                <input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="请输入密码"
                  className="w-full pl-12 h-12 bg-background/50 border border-white/10 rounded-full text-foreground placeholder:text-foreground/30 focus:border-primary transition-all px-4"
                  disabled={loading}
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full h-12 rounded-full bg-primary hover:bg-primary/90 text-white font-bold shadow-lg shadow-primary/30 transition-all active:scale-95 disabled:opacity-50"
            >
              {loading ? (
                <>
                  <Loader2 size={20} className="animate-spin inline mr-2" />
                  登录中...
                </>
              ) : (
                '登录'
              )}
            </button>
          </form>

          {/* 注册链接 */}
          <div className="mt-8 text-center">
            <p className="text-foreground/60 text-sm">
              还没有账号？{' '}
              <Link
                href="/register"
                className="text-primary font-semibold hover:underline transition-all"
              >
                立即注册
              </Link>
            </p>
          </div>
        </div>
      </motion.div>
    </div>
  );
}