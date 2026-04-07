import { NextRequest, NextResponse } from 'next/server';

const BACKEND_URL = 'http://localhost:8080';

async function proxyRequest(request: NextRequest, path: string) {
  const url = `${BACKEND_URL}/api/v1/${path}`;

  // Forward headers, excluding host and cookie (cookie rebuilt below)
  const headers = new Headers();
  request.headers.forEach((value, key) => {
    const k = key.toLowerCase();
    if (k !== 'host' && k !== 'cookie') {
      headers.set(key, value);
    }
  });

  // Rebuild cookie header from NextRequest.cookies — the reliable way
  // in Next.js App Router (request.headers.get('cookie') can be empty)
  const cookieParts: string[] = [];
  request.cookies.getAll().forEach((c) => {
    cookieParts.push(`${c.name}=${c.value}`);
  });
  if (cookieParts.length > 0) {
    headers.set('cookie', cookieParts.join('; '));
  }

  // Build fetch options
  const fetchOptions: RequestInit = {
    method: request.method,
    headers,
  };

  // Forward body for non-GET requests (use arrayBuffer to preserve binary data like file uploads)
  if (request.method !== 'GET' && request.method !== 'HEAD') {
    fetchOptions.body = await request.arrayBuffer();
  }

  // Make request to backend
  const backendResponse = await fetch(url, fetchOptions);

  // Read response as arrayBuffer to preserve binary responses too
  const responseBody = await backendResponse.arrayBuffer();
  const response = new NextResponse(responseBody, {
    status: backendResponse.status,
    statusText: backendResponse.statusText,
  });

  // Forward all response headers, especially Set-Cookie
  backendResponse.headers.forEach((value, key) => {
    // Handle multiple Set-Cookie headers
    if (key.toLowerCase() === 'set-cookie') {
      response.headers.append(key, value);
    } else {
      response.headers.set(key, value);
    }
  });

  return response;
}

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  return proxyRequest(request, path.join('/'));
}

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  return proxyRequest(request, path.join('/'));
}

export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  return proxyRequest(request, path.join('/'));
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  return proxyRequest(request, path.join('/'));
}

export async function PATCH(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  return proxyRequest(request, path.join('/'));
}
