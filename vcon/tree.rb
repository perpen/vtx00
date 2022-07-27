#!/usr/bin/env ruby
# FIXME
# - get rid of eq methods, use something else for tests.

=begin
createPanel(command: "bash", "~/src") => panelId
createPanel(content: "hello") => panelId
closePanel(1) => nil

layout(xxx) => nil
[
  {id: 0,
   pos: [0,0,80,25],
   content: "\n  hello\n",
   border: {
     style: {
       fg: blue,
       bg: black,
       bold: false,
     },
     title: {
       string: "main",
       style: {
         fg: blue,
         bg: black,
         bold: false,
       },
     },
     components: [
       {
         string: "✅",
         style: {
           fg: blue,
           bg: black,
           bold: false,
         },
       },
     ],
   },
]

state() =>

=end

require 'yaml'
require 'json'
require 'socket'
require 'thread'
require 'logger'

$log = Logger.new('tree.log')

class API

  def initialize queue
    @api_socket = "/tmp/vman-api.sock"
    @notif_socket = "/tmp/vman-notif.sock"
    @queue = queue
    Thread.new { listen }
  end

  def listen
    File.unlink(@notif_socket) if File.exists?(@notif_socket)
    socket = UNIXServer.new(@notif_socket)
    loop do
      conn = socket.accept
      reqJSON = conn.readline
      $log.debug "API.listen: received: #{reqJSON}"
      req = JSON.parse(reqJSON)
      @queue << req
    end
  end

  def send req
    @conn = UNIXSocket.new(@api_socket)
    json = JSON.generate(req)
    $log.debug "API.req: sending: #{json}"
    @conn.write json
    resp = JSON.parse(@conn.readline)
    $log.debug "API.req: reply: #{resp}"
    #@conn.close
    resp
  end
end

class Main

  attr_reader :client, :focus

  def start
    events_queue = Queue.new
    @api = API.new events_queue
    event_loop events_queue
  end

  def event_loop events_queue
    loop do
      loop do
        $log.debug "Main.event_loop"
        evt = events_queue.pop
        handle_event evt
      end
      $log.debug "Main.event_loop broke"
      sleep 1
    end
  end

  def handle_event evt
    evt_type = evt["EvtType"]
    handler = {
      "start" => method(:handle_start),
      "death" => method(:handle_death),
      "key" => method(:handle_key),
      "resize" => method(:handle_resize),
      }[evt_type]
    if handler.nil?
      $log.debug "Main.handle_event: unknown event type: #{evtType}"
      return
    end
    handler.call evt["EvtDetails"]
  end

  def handle_start args
    @w, @h = args

    @root = Node.new(VerticalStackLayout.new)
    @view = @root

    create_term
  end

  def handle_key args
    $log.debug("handle_key: focus=#{@focus}")
    @focus.handle_key(args, @focus) and return
    case args
    when 'c'.ord
      create_term
    when 'k'.ord
      focus_previous
    when 'j'.ord
      focus_next
    else
      $log.debug "Main.handle_key: unhandled key: #{args}"
    end
  end

  def handle_resize args
    @w, @h = args
    $log.debug "Main.handle_resize: #{@w}, #{@h}"
    push_updates @view, Position.new(0, 0, @w, @h)
  end

  def zoom node
    @view = node
    @focus = next_term node
  end

  def push_updates root, position=nil
    $log.debug "Main.push_updates: #{root}, #{position}"
    dirty_terms = root.arrange(position)
    states = dirty_terms.map {|term| term.to_api self}
    cmd = {
      :layoutCmd => {
        :focusId => @focus == nil ? nil : @focus.id,
        :panels => states,
      }
    }
    #$log.debug cmd
    @api.send cmd
  end

  def handle_death args
    panel_id = args
    node = node_by_id panel_id
    @focus = (next_term @focus or next_term @view)
    node.parent.remove_child {|n| n.id == panel_id }
    push_updates node.parent
  end

  def create_term
    term = Term.new(@api, ["/bin/bash"])
    @focus = term
    @view.add_child(term)
    push_updates @view, Position.new(0, 0, @w, @h)
   end

  def focus_previous
    updates = []
    old_focus = @focus
    @focus = @focus.nil? ? next_term(@view) : (previous_term(@focus) or next_term(@view))
    updates.push(old_focus.to_api self) if old_focus
    updates.push(@focus.to_api self) if @focus and not @focus.equal? old_focus
    cmd = {
      :layoutCmd => {
        :focusId => @focus == nil ? nil : @focus.id,
        :panels => updates,
      }
    }
    @api.send cmd
  end

  def focus_next
    updates = []
    old_focus = @focus
    @focus = @focus.nil? ? next_term(@view) : (next_term(@focus) or next_term(@view))
    updates.push(old_focus.to_api self) if old_focus
    updates.push(@focus.to_api self) if @focus and not @focus.equal? old_focus
    cmd = {
      :layoutCmd => {
        :focusId => @focus == nil ? nil : @focus.id,
        :panels => updates,
      }
    }
    @api.send cmd
  end

  def next_term node
    return next_node(@view, node) {|node| node.term?}
  end

  def previous_term node
    return previous_node(@view, node) {|node| node.term?}
  end

  def next_node root, current, &block
    nodes = all_nodes root
    use_next = false
    nodes.each do |node|
      return node if use_next and (not block_given? or block.call(node))
      use_next = true if node.equal? current
    end
    return nil
  end

  def previous_node root, current, &block
    nodes = all_nodes root
    previous = nil
    nodes.each do |node|
      return previous if node.equal? current
      previous = node if not block_given? or block.call(node)
    end
    return nil
  end

  def node_by_id id
    n = next_node(@root, @root) { |node| node.id == id }
    return n
  end

  def all_nodes root, &block
    nodes = []
    nodes.push(root)
    root.children.each do |child|
      nodes.concat(all_nodes(child))
    end
    return nodes
  end
end

# The provided key handling method passes the key to the layout, if layout does
# not handle it, forwards it to the parent node.
class Node

  attr_reader :children, :attributes, :layout
  attr_accessor :parent, :position

  def initialize(layout)
    @parent = nil
    @layout = layout
    @children = []
    @attributes = {}
    @position = nil
    #@position = Position.new(10, 10, 10, 10)
  end

  def term?
    false
  end

  def arrange position=nil
    $log.debug "Node.arrange: node: #{self} position: #{position}"
    position = (position or @position)
    return @layout.arrange(self, position)
  end

  def add_child(node, index: 0)
    node.parent = self
    @children.push(node)
  end

  def remove_child &block
    @children.delete_if &block
  end

  def add_attribute(name, value)
    @attributes[name] = value
  end

  def attribute(name)
    return @attributes[name]
  end

  def handle_key(k, term)
    $log.debug('Node.handle_key:', k, term)
    return (@layout.handle_key(k, term) or @parent.handle_key(k, term))
  end

  # Used by unit tests
  def ==(other)
    return (other.kind_of? Node and
            @parent == other.parent and
            @layout == other.layout and
            @attributes == other.attributes and
            @children == other.children)
  end
end

# A term is a panel showing a terminal.
class Term < Node

  attr_reader :cmd_argv, :id

  def initialize(client, cmd_argv)
    super(BaseLayout.new)
    @cmd_argv= cmd_argv
    req = {:CreatePanelCmd => {:Argv => cmd_argv, :Cwd => "."}}
    @id = client.send req
  end

  def term?
    true
  end

  def to_api ctx
    {:id => @id,
     :pos => [@position.x, @position.y, @position.w, @position.h],
     :border => {
       :style => ((equal? ctx.focus) ? Style::BORDER_FOCUSED : Style::BORDER_NORMAL).to_api,
       :title => {:string => "TITLE",
                  :style => Style::TITLE.to_api}
       }
     }
  end

  def to_s
    "Term(id: #{@id}, parent: #{@parent})"
  end

  # Used by unit tests
  def ==(other)
    return (other.instance_of? Term and
            @cmd_argv == other.cmd_argv)
  end
end

class Position
  attr_accessor :x, :y, :w, :h

  def initialize x, y, w, h
    @x, @y, @w, @h = x, y, w, h
  end

  def to_s
    "Position(x:#{@x} y:#{@y} w:#{@w} h:#{@h})"
  end

  def ==(other)
    return (other.instance_of? Position and
      @x == other.x and
      @y == other.y and
      @w == other.w and
      @h == other.h)
  end
end

class BaseLayout

  def arrange(node, position)
    if position != node.position
      node.position = position
      [node]
    else
      []
    end
  end

  def handle_key(k, term)
  end
end

# Layouts children nodes in an horizontal or vertical stack. Children are of equal size,
# except for those with a fixed-size attribute.
# TODO - cope with overflow
class StackLayout < BaseLayout

  def initialize(vertical)
    @vertical = vertical
  end

  def arrange(node, position)
    node.position = position
    updates = []
    x, y, w, h = position.x, position.y, position.w, position.h

    if node.children.length > 0
      # Compute children sizes
      size = @vertical ? h : w
      flexible_size = size
      flexible_children_count = node.children.length
      for child in node.children
        fixed_size = child.attribute("fixed-size")
        if fixed_size
          flexible_size -= fixed_size
          flexible_children_count -= 1
        end
      end
      flexible_child_size, leftover = flexible_size.divmod(flexible_children_count)

      offset = @vertical ? y : x
      node.children.each_with_index do |child, index|
        fixed_child_size = child.attribute("fixed-size")
        child_size = fixed_child_size || flexible_child_size
        child_size += index == 0 ? leftover : 0
        child_position = if @vertical
                           pos = Position.new(x, offset, w, child_size)
                         else
                           pos = Position.new(offset, y, child_size, h)
                         end
        child_updates = child.arrange(child_position)
        updates = updates.concat child_updates
        offset += child_size
      end
    end

    return updates
  end

  def handle_key args, term
    case args
    when '%'.ord
      node.parent.remove_child {|n| n.id == term.id }
      new_container = Node.new(HorizontalStackLayout.new)
      node.parent.add_child new_container
      new_container.add_child term
      new_term = create_term
      new_container.add_child new_term
      push_updates node.parent
    end
  end
end

class VerticalStackLayout < StackLayout
  def initialize
    super(true)
  end
end

class HorizontalStackLayout < StackLayout
  def initialize
    super(false)
  end
end

class Style

  def initialize fg, bg, bold
    @fg, @bg, @bold = fg, bg, bold
  end

  def to_api
    {:fg => @fg,
     :bg => @bg}
  end

  BORDER_NORMAL = Style.new "white", "black", false
  BORDER_FOCUSED = Style.new "red", "black", false
  TITLE = Style.new "white", "yellow", false
end

class BarComp
  def initialize string, style
    @string, @style = string, style
  end
end

class Border
  def initialize style, title, components
   @style, @title, @components = style, title, components
  end
end

class Panel
  def initialize id, pos, border, content
    @id, @pos, @border, @content = id, pos, border, content
  end
end

if ARGV.length > 0 and ARGV[0] == 'run'
  $log.level = Logger::DEBUG
  $log.info "Starting"
  begin
    main = Main.new
  rescue
    $log.error "oops"
    $log.error $!
  end

  main.start
end
