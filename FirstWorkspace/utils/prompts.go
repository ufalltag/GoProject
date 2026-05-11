package utils

const analyzePrompt = `You are a bookmark categorization assistant.

Analyze the following web page content.

Page content:
%s

%s

Your task:
1. Create a short, clear title for this bookmark (max 60 characters)
2. Suggest the most relevant category/folder name for this content
3. Rate your confidence from 0.0 to 1.0

Rules:
- Title must be specific and descriptive
- Folder name must be short (1-3 words)
- If existing folders are provided, prefer using one of them if it fits well
- If none fit, suggest a new folder name that best describes the content
- Respond ONLY with valid JSON, no extra text, no explanations

Examples:
{"title": "Swift Concurrency Guide", "folder": "iOS Development", "confidence": 0.9}
{"title": "PostgreSQL Performance Tips", "folder": "Databases", "confidence": 0.85}
{"title": "Figma Auto Layout Tutorial", "folder": "Design", "confidence": 0.8}

Required JSON format:
{
  "title": "short descriptive title",
  "folder": "category name",
  "confidence": 0.0 to 1.0
}`
