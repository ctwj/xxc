"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { motion } from 'framer-motion';
import { User, Mail, Lock, ArrowLeft, Loader2 } from 'lucide-react';
import { toast } from 'sonner';
import { api } from '@/lib/api';

export default function RegisterPage() {
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!username.trim() || !email.trim() || !password.trim()) {
      toast.error('请填写所有必填项');
      return;
    }

    // 验证用户名
    if (username.length < 3) {
      toast.error('用户名至少需要3个字符');
      return;
    }

    // 验证邮箱格式
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      toast.error('请输入有效的邮箱地址');
      return;
    }

    // 验证密码
    if (password.length < 6) {
      toast.error('密码至少需要6个字符');
      return;
    }

    // 确认密码
    if (password !== confirmPassword) {
      toast.error('两次输入的密码不一致');
      return;
    }

    setLoading(true);
    try {
      const response = await api.register(username, email, password);

      if (response.success) {
        toast.success('注册成功');
        router.push('/');
        router.refresh();
      } else {
        toast.error('注册失败：' + (response.error || '未知错误'));
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

      {/* 注册表单卡片 */}
      <motion.div
        initial={{ scale: 0.9, opacity: 0, y: 20 }}
        animate={{ scale: 1, opacity: 1, y: 0 }}
        transition={{ type: 'spring', damping: 20, stiffness: 300 }}
        className="relative z-10 w-full max-w-md mx-4"
      >
        <div className="glass-card rounded-[20px] p-10 shadow-2xl">
          {/* 标题 */}
          <div className="text-center mb-10">
            <h1 className="text-3xl font-bold mb-2 gradient-text">创建账号</h1>
            <p className="text-foreground/60 text-sm">注册以享受完整功能</p>
          </div>

          {/* 表单 */}
          <form onSubmit={handleSubmit} className="space-y-5">
            <div className="space-y-2">
              <label htmlFor="username" className="text-foreground/80 text-sm font-medium block">
                用户名
              </label>
              <div className="relative">
                <User size={18} className="absolute left-4 top-1/2 -translate-y-1/2 text-foreground/40" />
                <input
                  id="username"
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="请输入用户名"
                  className="w-full pl-12 h-12 bg-background/50 border border-white/10 rounded-full text-foreground placeholder:text-foreground/30 focus:border-primary transition-all px-4"
                  disabled={loading}
                />
              </div>
            </div>

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
                  placeholder="至少6个字符"
                  className="w-full pl-12 h-12 bg-background/50 border border-white/10 rounded-full text-foreground placeholder:text-foreground/30 focus:border-primary transition-all px-4"
                  disabled={loading}
                />
              </div>
            </div>

            <div className="space-y-2">
              <label htmlFor="confirmPassword" className="text-foreground/80 text-sm font-medium block">
                确认密码
              </label>
              <div className="relative">
                <Lock size={18} className="absolute left-4 top-1/2 -translate-y-1/2 text-foreground/40" />
                <input
                  id="confirmPassword"
                  type="password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  placeholder="请再次输入密码"
                  className="w-full pl-12 h-12 bg-background/50 border border-white/10 rounded-full text-foreground placeholder:text-foreground/30 focus:border-primary transition-all px-4"
                  disabled={loading}
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full h-12 rounded-full bg-primary hover:bg-primary/90 text-white font-bold shadow-lg shadow-primary/30 transition-all active:scale-95 disabled:opacity-50 mt-6"
            >
              {loading ? (
                <>
                  <Loader2 size={20} className="animate-spin inline mr-2" />
                  注册中...
                </>
              ) : (
                '注册'
              )}
            </button>
          </form>

          {/* 登录链接 */}
          <div className="mt-8 text-center">
            <p className="text-foreground/60 text-sm">
              已有账号？{' '}
              <Link
                href="/login"
                className="text-primary font-semibold hover:underline transition-all"
              >
                立即登录
              </Link>
            </p>
          </div>
        </div>
      </motion.div>
    </div>
  );
}
