"use client";

import { ScrollArea } from "@/components/ui/scroll-area";

export default function DataList({ items, onItemClick, selectedItem }: { items: string[], onItemClick: (item: string) => void, selectedItem: string }) {
  return (
    <ScrollArea className="h-full">
      <div className="space-y-1 p-2">
        {items.map((item, idx) => (
          <div
            key={idx}
            onClick={() => onItemClick(item)}
            className={`
              p-3 rounded-lg cursor-pointer
              transition-all duration-200
              hover:bg-gray-100
              ${selectedItem === item ? 'bg-gray-100 border-l-4 border-primary' : ''}
            `}
          >
            <h3 className="font-medium text-sm">{item}</h3>
          </div>
        ))}
      </div>
    </ScrollArea>
  );
}