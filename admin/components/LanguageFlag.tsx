"use client";

interface LanguageFlagProps {
  languageCode: string;
  className?: string;
}

const languageFlags: { [key: string]: string } = {
  'fr': '🇫🇷',
  'en': '🇬🇧', 
  'es': '🇪🇸',
  'de': '🇩🇪',
  'it': '🇮🇹',
  'pt': '🇵🇹',
  'zh': '🇨🇳',
  'ja': '🇯🇵',
  'ko': '🇰🇷',
  'ru': '🇷🇺',
  'ar': '🇸🇦',
  'hi': '🇮🇳',
  'tr': '🇹🇷',
  'pl': '🇵🇱',
  'nl': '🇳🇱',
  'sv': '🇸🇪',
  'da': '🇩🇰',
  'fi': '🇫🇮',
  'no': '🇳🇴',
  'cs': '🇨🇿',
  'sk': '🇸🇰',
  'hu': '🇭🇺',
  'ro': '🇷🇴',
  'bg': '🇧🇬',
  'hr': '🇭🇷',
  'sr': '🇷🇸',
  'sl': '🇸🇮',
  'et': '🇪🇪',
  'lv': '🇱🇻',
  'lt': '🇱🇹',
  'mt': '🇲🇹',
  'uk': '🇺🇦',
  'be': '🇧🇾',
  'mk': '🇲🇰',
  'sq': '🇦🇱',
  'ca': '🇪🇸',
  'eu': '🇪🇸',
  'gl': '🇪🇸',
  'cy': '🏴󠁧󠁢󠁷󠁬󠁳󠁿',
  'ga': '🇮🇪',
  'gd': '🏴󠁧󠁢󠁳󠁣󠁴󠁿',
  'br': '🇫🇷',
  'co': '🇫🇷',
  'oc': '🇫🇷',
  'lb': '🇱🇺',
  'rm': '🇨🇭',
  'fur': '🇮🇹',
  'sc': '🇮🇹',
  'vec': '🇮🇹',
  'lij': '🇮🇹',
  'pms': '🇮🇹',
  'lmo': '🇮🇹',
  'nap': '🇮🇹',
  'scn': '🇮🇹'
};

export default function LanguageFlag({ languageCode, className = "" }: LanguageFlagProps) {
  const flag = languageFlags[languageCode?.toLowerCase()] || '🏳️';
  
  return (
    <span 
      className={`inline-block ${className}`} 
      title={`Language: ${languageCode?.toUpperCase()}`}
    >
      {flag}
    </span>
  );
}