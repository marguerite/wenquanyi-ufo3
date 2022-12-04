#!/usr/bin/env ruby

class Array
  def group(n)
    arr = []
    0.step(size, n) do |i|
      arr << self[i, n]
    end
    arr
  end
end

ARGV.each do |d|
  arr = []
  arr1 = []
  Dir.foreach(File.join(d, 'glyphs')) do |f|
    next if ['.', '..'].include?(f)

    f1 = File.join(d, 'glyphs', f)
    b = false
    File.open(f1, 'r').each do |line|
      if /^.*<unicode hex.*$/ =~ line
        b = true
        break
      end
    end

    if !b && !File.basename(f1).end_with?('.plist')
      arr << '\\' + File.basename(f1).gsub('.glif', '').gsub('_', '')
      if (File.basename(f1).start_with?('H_') || File.basename(f1).start_with?('u')) && !File.basename(f1).end_with?('.mono.glif')
        arr1 << f1
      end
    end
  end

  arr.group(10).each { |a| puts a.join(' ') }

  f2 = File.open('features.fea', 'r')
  str = f2.read
  arr.each do |a|
    str = str.gsub(a, '')
  end
  f2.close
  File.open('features.fea', 'w') { |f| f.write str }

  arr1.each do |a|
    content = ''
    File.open(a, 'r').each do |line|
      if /^.*<advance.*$/ =~ line
        content += line
        hex = File.basename(a).gsub('.glif', '').gsub('_', '').gsub('uni', '')
                  .gsub('u', '').gsub('H', '')
        content += "\s\s<unicode hex=\"" + hex + "\"/>\n"
      else
        content += line
      end
    end
    File.open(a, 'w') { |f| f.write content }
  end
end
