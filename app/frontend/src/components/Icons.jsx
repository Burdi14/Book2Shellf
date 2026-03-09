import React from 'react';

const baseProps = {
  fill: 'none',
  stroke: 'currentColor',
  strokeLinecap: 'round',
  strokeLinejoin: 'round'
};

export function DownloadIcon({ size = 18, strokeWidth = 1.8, ...props }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" {...baseProps} strokeWidth={strokeWidth} {...props}>
      <path d="M12 3v12" />
      <path d="m7 11 5 5 5-5" />
      <path d="M5 19h14" />
    </svg>
  );
}

export function BookIcon({ size = 18, strokeWidth = 1.8, ...props }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" {...baseProps} strokeWidth={strokeWidth} {...props}>
      <path d="M5 4h9a4 4 0 0 1 4 4v11H9a4 4 0 0 0-4 4V4z" />
      <path d="M9 4v15" />
    </svg>
  );
}

export function FolderIcon({ size = 18, strokeWidth = 1.8, ...props }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" {...baseProps} strokeWidth={strokeWidth} {...props}>
      <path d="M3 6h6l2 3h10v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V6z" />
    </svg>
  );
}

export function DriveIcon({ size = 18, strokeWidth = 1.8, ...props }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" {...baseProps} strokeWidth={strokeWidth} {...props}>
      <rect x="3" y="7" width="18" height="12" rx="2" />
      <path d="M6 11h.01M10 11h.01M14 17h4" />
    </svg>
  );
}

export function PdfIcon({ size = 18, strokeWidth = 1.8, ...props }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" {...baseProps} strokeWidth={strokeWidth} {...props}>
      <path d="M7 2h7l5 5v15H7z" />
      <path d="M14 2v5h5" />
      <path d="M9 13h1.5a1.5 1.5 0 0 0 0-3H9v6" />
      <path d="M13 10v6" />
      <path d="M13 13h2" />
      <path d="M17 10h2a1 1 0 0 1 1 1v1a2 2 0 0 1-2 2h-1v-4Z" />
    </svg>
  );
}

export function ImageIcon({ size = 18, strokeWidth = 1.8, ...props }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" {...baseProps} strokeWidth={strokeWidth} {...props}>
      <rect x="3" y="5" width="18" height="14" rx="2" />
      <path d="m7 13 3-3 4 4 2-2 3 3" />
      <circle cx="9" cy="9" r="1.5" />
    </svg>
  );
}

export function EditIcon({ size = 16, strokeWidth = 1.8, ...props }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" {...baseProps} strokeWidth={strokeWidth} {...props}>
      <path d="M12 20h9" />
      <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4Z" />
    </svg>
  );
}

export function TrashIcon({ size = 16, strokeWidth = 1.8, ...props }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" {...baseProps} strokeWidth={strokeWidth} {...props}>
      <path d="M4 7h16" />
      <path d="M10 11v6" />
      <path d="M14 11v6" />
      <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2l1-12" />
      <path d="M9 7V4h6v3" />
    </svg>
  );
}

export function FilesIcon({ size = 18, strokeWidth = 1.8, ...props }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" {...baseProps} strokeWidth={strokeWidth} {...props}>
      <path d="M8 4h9l3 3v11a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2Z" />
      <path d="M14 4v4h4" />
      <path d="M6 8H5a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h9" />
    </svg>
  );
}

export function DownloadsIcon({ size = 16, strokeWidth = 1.8, ...props }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" {...baseProps} strokeWidth={strokeWidth} {...props}>
      <path d="M12 3v11" />
      <path d="m7 11 5 5 5-5" />
      <path d="M5 19h14" />
    </svg>
  );
}

export function SizeIcon({ size = 16, strokeWidth = 1.8, ...props }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" {...baseProps} strokeWidth={strokeWidth} {...props}>
      <rect x="3" y="5" width="18" height="14" rx="3" />
      <path d="M7 9h10" />
      <path d="M7 13h5" />
    </svg>
  );
}

export function CalendarIcon({ size = 16, strokeWidth = 1.8, ...props }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" {...baseProps} strokeWidth={strokeWidth} {...props}>
      <rect x="3" y="4" width="18" height="18" rx="2" />
      <path d="M7 2v4M17 2v4" />
      <path d="M3 10h18" />
      <path d="M8 14h2M12 14h2M16 14h2" />
    </svg>
  );
}
